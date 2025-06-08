const puppeteer = require("puppeteer");
const httpServer = require("http-server");
const path = require("path");
const assert = require("assert");

const serverHost = "127.0.0.1";
const serverPort = 8123;
const serverUrl = `http://${serverHost}:${serverPort}`;

const initialOutput = "---\nexample:\n  - first\n  - second\n";
const initialHash =
  "t=LQhQFMA8EMFsAcA24BcACA3h4aB0A1aRAV3AGdcBLAF3FjLQB81qB7ATTkSbQDtLeAE3C9qaAExoAvlNBA&v=LQhQEsBcFMFsGcBcACA2gIgGbgE70ugDTLrzQDGA9gHYAm6AuqEA";
const changeValuesContentTo = 'items: ["hello", "from", "puppeteer"]';
const expectedOutput = "---\nexample:\n  - hello\n  - from\n  - puppeteer\n";
const expectedHash =
  "t=LQhQFMA8EMFsAcA24BcACA3h4aB0A1aRAV3AGdcBLAF3FjLQB81qB7ATTkSbQDtLeAE3C9qaAExoAvlNBA&v=JYFwpgtgzgXABAbQEQAswBt0HskBo5IBmATlhHgQA4CullY4YxSAukA";

async function run() {
  console.log("[browser-test] starting server");
  const server = httpServer.createServer({
    root: path.join(__dirname, ".."),
    cache: -1,
  });
  await new Promise((resolve, reject) => {
    server.listen(serverPort, serverHost, (err) => {
      if (err) {
        return reject(err);
      }
      resolve();
    });
  });
  console.log(`[browser-test] server started at ${serverUrl}`);

  console.log("[browser-test] starting browser");
  const browser = await puppeteer.launch({ headless: true, devtools: false, args: ['--no-sandbox', '--disable-setuid-sandbox'] });
  const page = await browser.newPage();

  console.debug("Blank URL: " + page.url());

  try {
    console.group("Navigating to the server URL...");
    console.debug("Initial URL: " + serverUrl);
    let currentUrl = serverUrl;
    while (true) {
      console.debug("Current URL: " + currentUrl);
      await page.goto(currentUrl, { waitUntil: 'networkidle0' });
      const newUrl = page.url();
      if (newUrl === currentUrl) {
        console.debug("Navigated to: " + currentUrl);
        break;
      } else {
        currentUrl = newUrl;
      }
    }
    console.groupEnd();

    const getPageUrl = async () =>
      page.evaluate(() => window.location.toString());
    const setValuesContent = async (value) =>
      page.$eval(
        ".textarea--values",
        (el, value) => el.CodeMirror.setValue(value),
        value
      );
    const getOutput = async () =>
      page.$eval(".textarea--output", (el) => el.CodeMirror.getValue());

    console.log("[browser-test] running test");
    assert.equal(await getOutput(), initialOutput, "should see initial output");
    assert.equal(
      await getPageUrl(),
      `${serverUrl}/#${initialHash}`,
      "should see initial URL hash"
    );
    await setValuesContent(changeValuesContentTo);
    assert.equal(await getOutput(), expectedOutput, "should see updated output");
    assert.equal(
      await getPageUrl(),
      `${serverUrl}/#${expectedHash}`,
      "should see updated URL hash"
    );
    console.log("[browser-test] ran test");
    console.log("[browser-test] PASS");
  } catch (error) {
    console.error("Error:", error);
  } finally {
    await browser.close();
    console.debug("[browser-test] closed browser");
    server.close();
    console.debug("[browser-test] closed server");
  }
}

run();
