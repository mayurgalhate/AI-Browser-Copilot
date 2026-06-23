declare const chrome: any;
import { extractInteractiveElements, extractContent } from './extractor';
import { clickElement, fillInput, scroll, browserHistory, getCurrentUrl } from './interactor';

chrome.runtime.onMessage.addListener((message: any, _sender: any, sendResponse: any) => {
  if (message.type === "EXECUTE_ACTION") {
    handleAction(message).then(res => sendResponse({ result: res })).catch(err => sendResponse({ result: `Error: ${err.message}` }));
    return true; // Keep message channel open for async response
  }
});

async function handleAction(msg: any): Promise<string> {
  const { action, target, value } = msg;

  switch (action) {
    case "extract_elements":
      return extractInteractiveElements();
    case "click":
      return clickElement(target);
    case "fill_input":
      let id = target;
      let val = value;
      try {
        const parsed = JSON.parse(target);
        if (parsed.id) id = parsed.id;
        if (parsed.value) val = parsed.value;
      } catch (e) {}
      return fillInput(id, val);
    case "scroll":
      return scroll(target || "down");
    case "browser_history":
      return browserHistory(target || "back");
    case "get_current_url":
      return getCurrentUrl();
    case "extract_content":
      return extractContent();
    default:
      throw new Error(`Unknown action: ${action}`);
  }
}
