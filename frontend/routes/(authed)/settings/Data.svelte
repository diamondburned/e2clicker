<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import QRCode from "$lib/components/QRCode.svelte";
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";

  import * as api from "$lib/api.svelte";
  import { user } from "$lib/api.svelte";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  type Status<DataT extends object = {}> =
    | { status: undefined }
    | { status: "loading" }
    | {
        status: "error";
        error: any;
      }
    | ({
        status: "success";
      } & DataT);

  type SupportedMIME = "text/csv" | "application/json";

  let importStatus = $state<
    Status<{
      succeeded: number;
      records: number;
    }>
  >({ status: undefined });

  let exportStatus = $state<Status>({ status: undefined });

  function filenameToMIME(filename: string): SupportedMIME {
    const ext = filename.split(".").pop();
    return (
      {
        csv: "text/csv",
        json: "application/json",
      } as Record<string, SupportedMIME>
    )[ext ?? ""];
  }

  async function importFrom(file: File) {
    try {
      importStatus = { status: "loading" };

      const mime = filenameToMIME(file.name);
      if (!mime) {
        throw new Error("Invalid file type.");
      }

      // Use type File as parameter instead of a string.
      const resp = await api.fetch("/dosage/import", {
        method: "POST",
        body: file,
        headers: {
          "Content-Type": mime,
        },
      });

      const body = await (resp.json() as ReturnType<typeof api.importDosageHistory>);
      if (body.error) {
        throw new Error(body.error.message, { cause: body.error });
      }

      importStatus = { status: "success", ...body };
    } catch (err) {
      importStatus = { status: "error", error: err };
    }
  }

  async function exportTo(format: SupportedMIME) {
    try {
      exportStatus = { status: "loading" };

      const resp = await api.fetch("/dosage/export", {
        method: "GET",
        headers: {
          Accept: format,
        },
      });

      const filename =
        resp.headers.get("Content-Disposition")?.split("filename=")[1] ?? "dosage.txt";
      const object = URL.createObjectURL(await resp.blob());

      const a = document.createElement("a");
      a.download = filename;
      a.href = object;
      a.click();

      exportStatus = { status: "success" };
    } catch (err) {
      exportStatus = { status: "error", error: err };
    }
  }
</script>

<PreferenceGroup name="Data">
  {#snippet description()}
    Manage your data.
  {/snippet}

  <PreferenceItem name="Import" for="import-history-button">
    {#snippet description()}
      Import your dose history from CSV or JSON. You may do this 3 times maximum every 15 minutes.
      Doses that match an existing timestamp will be ignored.

      <span class="block mt-2 empty:hidden">
        {#if importStatus?.status == "success"}
          <span class="status-success text-green">
            <b>{importStatus.succeeded}</b> records imported out of
            <b>{importStatus.records}</b>.
          </span>
        {:else if importStatus?.status == "error"}
          <ErrorBox error={importStatus.error} tiny="inline" />
        {/if}
      </span>
    {/snippet}

    <button
      id="import-history-button"
      class="outline ml-2"
      onclick={() => document.getElementById("import-history")?.click()}
      disabled={importStatus?.status == "loading"}
    >
      <Icon name="upload" />
      Import
    </button>

    <input
      id="import-history"
      type="file"
      class="hidden"
      disabled={importStatus?.status == "loading"}
      onchange={(ev) => {
        const file = ev.currentTarget.files?.[0];
        if (file) {
          importFrom(file);
        }
      }}
    />
  </PreferenceItem>

  <PreferenceItem name="Export" for="">
    {#snippet description()}
      Export your dose history as CSV or JSON. You may do this 3 times maximum every 15 minutes.

      <span class="block mt-2 empty:hidden">
        {#if exportStatus?.status == "error"}
          <ErrorBox error={exportStatus.error} tiny="inline" />
        {/if}
      </span>
    {/snippet}

    <button
      class="outline"
      onclick={() => exportTo("application/json")}
      disabled={exportStatus?.status == "loading"}
    >
      JSON
      <Icon name="download" />
    </button>
    <button
      class="outline ml-2"
      onclick={() => exportTo("text/csv")}
      disabled={exportStatus?.status == "loading"}
    >
      CSV
      <Icon name="download" />
    </button>
  </PreferenceItem>
</PreferenceGroup>
