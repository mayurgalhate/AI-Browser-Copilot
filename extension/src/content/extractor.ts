export function extractInteractiveElements(): string {
  const elements = document.querySelectorAll('button, a, input, select, textarea');
  const result: any[] = [];
  
  elements.forEach((el, index) => {
    // Make sure the element is visible
    const rect = el.getBoundingClientRect();
    if (rect.width > 0 && rect.height > 0) {
      // Assign a stable ID if it doesn't have one
      let stableId = el.id || el.getAttribute('data-copilot-id');
      if (!stableId) {
        stableId = `copilot-el-${index}`;
        el.setAttribute('data-copilot-id', stableId);
      }
      
      result.push({
        id: stableId,
        tag: el.tagName.toLowerCase(),
        text: (el as HTMLElement).innerText || (el as HTMLInputElement).value || el.getAttribute('aria-label') || 'unnamed',
        type: el.getAttribute('type') || undefined
      });
    }
  });

  return JSON.stringify(result);
}

export function extractContent(): string {
  return document.body.innerText.substring(0, 5000); // return up to 5k chars
}
