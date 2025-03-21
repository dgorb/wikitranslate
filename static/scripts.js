let languages = {};

async function populateLanguageDropdowns() {
  try {
    const response = await fetch("/languages");
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
}

async function performTranslation() {
  const inputLang = document.getElementById("inputLang").value;
  const outputLang = document.getElementById("outputLang").value;
  const input = document.getElementById("input").value;
  const resultsDiv = document.getElementById("results");
  const translationDiv = document.getElementById("translation");
  const summaryDiv = document.getElementById("summary");
  resultsDiv.classList.remove("hidden");

  try {
    const response = await fetch(
      `/translate?inputLang=${inputLang}&outputLang=${outputLang}&input=${encodeURIComponent(input)}`,
    );
    const data = await response.json();
    if (response.ok) {
      translationDiv.innerHTML = `${input.charAt(0).toUpperCase() + input.slice(1)} (${inputLang})  →  ${data.translation} (${outputLang})`;
      summaryDiv.innerHTML = `${data.summary}`;
    } else {
      resultsDiv.innerHTML = `<div class="text-red-600">Translation failed</div>`;
    }
  } catch (error) {
    console.log(error);
    resultsDiv.innerHTML = `<div class="text-red-600">Translation not found</div>`;
  }

  localStorage.setItem("lastInputLang", inputLang);
  localStorage.setItem("lastOutputLang", outputLang);
}

document.addEventListener("DOMContentLoaded", populateLanguageDropdowns);
