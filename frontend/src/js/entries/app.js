import "../../css/app.css";
import "vite/modulepreload-polyfill";
import "htmx.org";
import u from "umbrellajs";
import { writeTextToClipboard } from "~/lib/clipboard";

u(document)
  .on("click", ".js-copy", (ev) => {
    const target = ev.currentTarget;
    writeTextToClipboard(target.dataset.copyValue);
    const checkmark = u(target).parent().find(".js-checkmark");
    checkmark.removeClass("opacity-0");
    setTimeout(() => {
      checkmark.addClass("opacity-0");
    }, 2000);
  })
  .on("click", ".js-drawer-open", () => {
    u("#drawer").removeClass("translate-x-full", "opacity-0");
  })
  .on("click", ".js-drawer-close", () => {
    u("#drawer").addClass("translate-x-full");
  })
  .on("click", "#archiveButton", () => {
    u("#archiveDialog").first().showModal();
  })
  .on("click", "#recoverButton", () => {
    u("#recoverDialog").first().showModal();
  })
  .on("click", ".js-close-dialog", (ev) => {
    u(ev.currentTarget).closest("dialog").first().close();
  });

u(".js-config-length-toggle").on("change", (ev) => {
  const enabled = ev.currentTarget.checked;
  const panel = u(ev.currentTarget).closest(".js-generate-configuration");

  if (enabled) {
    panel.find(".js-config-length-opacity").removeClass("opacity-40");
  } else {
    panel.find(".js-config-length-opacity").addClass("opacity-40");
  }

  panel.find(".js-length-linked").each((input) => {
    input.disabled = !enabled;
  });
});

u(".js-length-range-slider").on("input", (ev) => {
  const valueEl = u(ev.currentTarget)
    .closest(".js-length-range-container")
    .find(".js-length-range-value")
    .first();
  if (valueEl) {
    valueEl.textContent = ev.currentTarget.value;
  }
});

u(".js-filter-length-toggle").on("change", (ev) => {
  const enabled = ev.currentTarget.checked;
  const form = u(ev.currentTarget).closest("form");

  if (enabled) {
    form.find(".js-filter-length-controls").removeClass("opacity-40");
  } else {
    form.find(".js-filter-length-controls").addClass("opacity-40");
  }

  form.find(".js-filter-length-linked").each((input) => {
    input.disabled = !enabled;
  });
});

u(".js-filter-length-range-slider").on("input", (ev) => {
  const valueEl = u(ev.currentTarget)
    .closest(".js-filter-length-range-container")
    .find(".js-filter-length-range-value")
    .first();
  if (valueEl) {
    valueEl.textContent = ev.currentTarget.value;
  }
});
