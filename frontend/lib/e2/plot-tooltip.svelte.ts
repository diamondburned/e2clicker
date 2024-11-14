import type * as charts from "lightweight-charts";
import { precision, units } from "./plot.svelte";
import { DateTime } from "luxon";

export type PlotTooltipData = {
  left: number;
  top: number;
  time: DateTime;
  data: { name: string; value: string }[];
  visible?: boolean;
};

export type AcceptablePlotTooltipData = {
  name: string;
  value: number; // NaN acceptable
  unit?: string;
}[];

function ptInElement(pt: { x: number; y: number }, elem: HTMLElement) {
  return pt.x > 0 && pt.y > 0 && pt.x < elem.clientWidth && pt.y < elem.clientHeight;
}

function fixed(x?: number) {
  return x && x.toFixed(precision());
}

const tooltipMargin = 8;

export function renderPlotTooltip(
  oldData: PlotTooltipData | null,
  ev: charts.MouseEventParams<charts.Time>,
  plotDiv: HTMLDivElement | null,
  plotTooltipDiv: HTMLDivElement | null,
  acceptedData: (AcceptablePlotTooltipData[number] | undefined)[],
): PlotTooltipData | null {
  if (!plotDiv || !plotTooltipDiv || !ev.point || !ev.time || !ptInElement(ev.point, plotDiv)) {
    if (!oldData) {
      return null;
    }
    return {
      ...oldData,
      visible: false,
    };
  }

  const { clientWidth: tooltipWidth, clientHeight: tooltipHeight } = plotTooltipDiv;
  const y = ev.point.y;

  let left = ev.point.x + tooltipMargin * 2 + tooltipWidth / 2;
  if (left > plotDiv.clientWidth - tooltipWidth) {
    left = ev.point.x - tooltipMargin - tooltipWidth / 2;
  }

  let top = y + tooltipMargin;
  if (top > plotDiv.clientHeight - tooltipHeight) {
    top = y - tooltipHeight;
  }

  const time = DateTime.fromSeconds(ev.time as charts.UTCTimestamp);
  const data = acceptedData
    .filter((d) => !!d)
    .filter((d) => !isNaN(d.value))
    .map((d) => ({
      name: d.name,
      value: `${fixed(d.value)} ${d.unit ?? units()}`,
    }));

  return {
    left,
    top,
    time,
    data,
    visible: data.length > 0,
  };
}

export function lineDataFromSeriesData(
  seriesData: charts.MouseEventParams<charts.Time>["seriesData"],
  series: charts.ISeriesApi<"Line", charts.UTCTimestamp>,
): charts.LineData<charts.UTCTimestamp> {
  return seriesData.get(
    series as charts.ISeriesApi<"Line">,
  ) as charts.LineData<charts.UTCTimestamp>;
}
