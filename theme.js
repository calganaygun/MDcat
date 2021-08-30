function switchTheme(isLight) {
    const markdownBodyElement = document.getElementById("markdown-body");
    const themeIconElement = document.getElementById("theme-icon");
    const lightHL = document.getElementById("light-hl");
    const darkHL = document.getElementById("dark-hl");

    const LS = window.localStorage;
    const LSisLight = LS.getItem("MDcatisLight") === "true";

    isLight = isLight ?? !LSisLight;

    const colorMode = isLight ? "light" : "dark";
    const colorIcon = isLight ? "sun" : "moon";

    markdownBodyElement.setAttribute("data-color-mode", colorMode);
    markdownBodyElement.setAttribute("data-dark-theme", colorMode);
    lightHL.rel =  isLight ? "stylesheet" : "stylesheet alternate";
    darkHL.rel = isLight ? "stylesheet alternate" : "stylesheet";

    themeIconElement.setAttribute(
      "data-icon",
      "octicon:" + colorIcon + "-16"
    );
    LS.setItem("MDcatisLight", isLight);
  }

  const themeButton = document.getElementById("theme-button");
  const LS = window.localStorage;
  const LSisLight = LS.getItem("MDcatisLight") === "true";

  themeButton.addEventListener("click", () => {
    switchTheme();
  });

  LS.getItem("MDcatisLight") === null
    ? switchTheme(
        !(
          window.matchMedia &&
          window.matchMedia("(prefers-color-scheme: dark)").matches
        )
      )
    : switchTheme(LSisLight);