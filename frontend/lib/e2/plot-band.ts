import type * as charts from "lightweight-charts";

type BandPoint = {
  x: number;
  yTop: number;
  yBot: number;
};

type DrawTarget = Parameters<charts.ISeriesPrimitivePaneRenderer["draw"]>[0];

class UncertaintyBandRenderer implements charts.ISeriesPrimitivePaneRenderer {
  private readonly points: readonly BandPoint[];
  private readonly color: string;

  constructor(points: readonly BandPoint[], color: string) {
    this.points = points;
    this.color = color;
  }

  draw(target: DrawTarget): void {
    const points = this.points;
    if (points.length < 2) {
      return;
    }

    target.useMediaCoordinateSpace(({ context: ctx }) => {
      ctx.save();
      ctx.fillStyle = this.color;
      ctx.beginPath();
      ctx.moveTo(points[0].x, points[0].yTop);
      for (let i = 1; i < points.length; i++) {
        ctx.lineTo(points[i].x, points[i].yTop);
      }
      for (let i = points.length - 1; i >= 0; i--) {
        ctx.lineTo(points[i].x, points[i].yBot);
      }
      ctx.closePath();
      ctx.fill();
      ctx.restore();
    });
  }
}

class UncertaintyBandPaneView implements charts.ISeriesPrimitivePaneView {
  private readonly primitive: UncertaintyBandPrimitive;

  constructor(primitive: UncertaintyBandPrimitive) {
    this.primitive = primitive;
  }

  zOrder(): charts.SeriesPrimitivePaneViewZOrder {
    return "bottom";
  }

  renderer(): charts.ISeriesPrimitivePaneRenderer | null {
    if (!this.primitive.visible || this.primitive.points.length < 2) {
      return null;
    }
    return new UncertaintyBandRenderer(this.primitive.points, this.primitive.color);
  }
}

export class UncertaintyBandPrimitive implements charts.ISeriesPrimitive<charts.UTCTimestamp> {
  private chart: charts.IChartApiBase<charts.UTCTimestamp> | null = null;
  private series: charts.ISeriesApi<charts.SeriesType, charts.UTCTimestamp> | null = null;
  private requestUpdate: (() => void) | null = null;

  private lower: readonly charts.LineData<charts.UTCTimestamp>[] = [];
  private upper: readonly charts.LineData<charts.UTCTimestamp>[] = [];
  private readonly paneViewsList: readonly charts.ISeriesPrimitivePaneView[];

  color = "transparent";
  visible = true;
  points: BandPoint[] = [];

  constructor(options?: { color?: string; visible?: boolean }) {
    if (options?.color !== undefined) {
      this.color = options.color;
    }
    if (options?.visible !== undefined) {
      this.visible = options.visible;
    }
    this.paneViewsList = [new UncertaintyBandPaneView(this)];
  }

  attached(param: charts.SeriesAttachedParameter<charts.UTCTimestamp, charts.SeriesType>): void {
    this.chart = param.chart;
    this.series = param.series;
    this.requestUpdate = param.requestUpdate;
    this.updateAllViews();
  }

  detached(): void {
    this.chart = null;
    this.series = null;
    this.requestUpdate = null;
    this.points = [];
  }

  setData(
    lower: readonly charts.LineData<charts.UTCTimestamp>[],
    upper: readonly charts.LineData<charts.UTCTimestamp>[],
  ): void {
    this.lower = lower;
    this.upper = upper;
    this.requestUpdate?.();
  }

  setOptions(options: { color?: string; visible?: boolean }): void {
    if (options.color !== undefined) {
      this.color = options.color;
    }
    if (options.visible !== undefined) {
      this.visible = options.visible;
    }
    this.requestUpdate?.();
  }

  updateAllViews(): void {
    if (
      !this.visible ||
      !this.chart ||
      !this.series ||
      this.lower.length === 0 ||
      this.upper.length === 0
    ) {
      this.points = [];
      return;
    }

    const timeScale = this.chart.timeScale();
    const series = this.series;
    const pts: BandPoint[] = [];

    let i = 0;
    let j = 0;
    while (i < this.lower.length && j < this.upper.length) {
      const lo = this.lower[i];
      const hi = this.upper[j];
      if (lo.time < hi.time) {
        i++;
      } else if (hi.time < lo.time) {
        j++;
      } else {
        const x = timeScale.timeToCoordinate(lo.time);
        const yTop = series.priceToCoordinate(hi.value);
        const yBot = series.priceToCoordinate(lo.value);
        if (x !== null && yTop !== null && yBot !== null) {
          pts.push({ x, yTop, yBot });
        }
        i++;
        j++;
      }
    }

    this.points = pts;
  }

  paneViews(): readonly charts.ISeriesPrimitivePaneView[] {
    return this.paneViewsList;
  }
}
