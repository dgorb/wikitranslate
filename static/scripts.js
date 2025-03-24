let languages = {};

async function populateLanguageDropdowns() {
  try {
    const response = await fetch("/api/languages");
    if (!response.ok) {
      throw new Error("Failed to fetch languages");
    }
    languages = await response.json();

    const inputLang = document.getElementById("inputLang");
    const outputLang = document.getElementById("outputLang");

    inputLang.innerHTML = "";
    outputLang.innerHTML = "";

    const sortedLanguages = Object.entries(languages).sort((a, b) =>
      a[1].localeCompare(b[1]),
    );

    for (const [code, name] of sortedLanguages) {
      const option1 = new Option(name, code);
      const option2 = new Option(name, code);

      inputLang.add(option1);
      outputLang.add(option2);

      inputLang.value =
        code === "en"
          ? code
          : localStorage.getItem("lastInputLang") || inputLang.value;
      outputLang.value =
        code === "nb"
          ? code
          : localStorage.getItem("lastOutputLang") || outputLang.value;
    }
  } catch (error) {
    console.error("Error loading languages:", error);
    const resultDiv = document.getElementById("result");
    resultDiv.classList.remove("hidden");
    resultDiv.innerHTML = `<div class="text-red-600">Error loading languages: ${error.message}</div>`;
  }

  document.getElementById("input").focus();
}

async function swapLanguage() {
  const inputLangSelect = document.getElementById("inputLang");
  const outputLangSelect = document.getElementById("outputLang");
  const inputEl = document.getElementById("input");

  const input = document
    .getElementById("translation")
    .getAttribute("data-translation");

  [inputLangSelect.value, outputLangSelect.value] = [
    outputLangSelect.value,
    inputLangSelect.value,
  ];
  inputEl.value = input;

  performTranslation();
}

async function performTranslation() {
  const inputLang = document.getElementById("inputLang").value;
  const outputLang = document.getElementById("outputLang").value;
  const input = document.getElementById("input");
  const inputVal = input.value;
  const resultsDiv = document.getElementById("results");
  const translationDiv = document.getElementById("translation");
  const summaryDiv = document.getElementById("summary");
  const swapButton = document.getElementById("swapSummaryButton");

  resultsDiv.classList.remove("hidden");
  swapButton.classList.remove("hidden");

  try {
    const response = await fetch(
      `/api/translate?inputLang=${inputLang}&outputLang=${outputLang}&input=${encodeURIComponent(inputVal)}`,
    );
    const data = await response.json();
    if (response.ok) {
      translationDiv.setAttribute("data-translation", data.translation);
      translationDiv.innerHTML = `${inputVal.charAt(0).toUpperCase() + inputVal.slice(1)} (${inputLang})  →  ${data.translation} (${outputLang})`;

      summaryDiv.innerHTML = data.summary;
      currentSummaryView = "output";
      swapButton.textContent = "🔁";
    } else {
      translationDiv.innerHTML = "Error translating";
      summaryDiv.innerHTML = "";
      swapButton.classList.add("hidden");
    }
  } catch (error) {
    console.log(error);
    translationDiv.innerHTML = "Translation not found";
    summaryDiv.innerHTML = "";
    swapButton.classList.add("hidden");
  }

  if (inputVal) {
    const newUrl = `/translate?inputLang=${inputLang}&outputLang=${outputLang}&input=${encodeURIComponent(inputVal)}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
  }

  localStorage.setItem("lastInputLang", inputLang);
  localStorage.setItem("lastOutputLang", outputLang);
}

function checkUrlForTranslation() {
  const urlParams = new URLSearchParams(window.location.search);
  const inputLang = urlParams.get("inputLang");
  const outputLang = urlParams.get("outputLang");
  const input = urlParams.get("input");

  if (inputLang && outputLang && input) {
    setTimeout(() => {
      document.getElementById("inputLang").value = inputLang;
      document.getElementById("outputLang").value = outputLang;
      document.getElementById("input").value = input;

      performTranslation();

      const newUrl = `/translate?inputLang=${inputLang}&outputLang=${outputLang}&input=${encodeURIComponent(input)}`;
      window.history.pushState({ path: newUrl }, "", newUrl);
    }, 50);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  populateLanguageDropdowns();
  checkUrlForTranslation();
});
