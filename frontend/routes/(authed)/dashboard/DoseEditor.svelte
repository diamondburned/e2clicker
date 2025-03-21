<script lang="ts">
  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api";

  import Dialog from "$lib/components/Dialog.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import InputQuantity from "$lib/components/inputs/InputQuantity.svelte";

  let {
    open = $bindable(false),
    dose,
  }: {
    open: boolean;
    dose?: e2.Dose;
  } = $props();

  async function saveDose(dose: e2.Dose) {
    console.log("Saving dose", dose);
  }
</script>

{#if open && dose}
  <Dialog bind:open dismissible>
    <header>
      <h2>Edit Dose</h2>
    </header>
    <form>
      <p>Editing dose: {JSON.stringify(dose)}</p>
      <label>
        <span>Dose Time</span>
        <input type="datetime-local" value={dose.takenAt.toISO()} />
      </label>
      <label>
        <span>Dose Amount</span>
        <InputQuantity
          unit={dose.deliveryMethod.units}
          initial={dose.dose}
          onchange={(qty) => {
            if (qty) {
              dose!.dose = qty;
            }
          }}
        />
      </label>
    </form>
    <footer>
      <button aria-label="Cancel" onclick={() => (dose = undefined)}>
        Cancel <Icon name="close" />
      </button>
      <button aria-label="Save" onclick={() => saveDose(dose!)}>
        Save <Icon name="save" />
      </button>
    </footer>
  </Dialog>
{/if}
