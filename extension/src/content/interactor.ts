export function clickElement(id: string): string {
  const el = document.getElementById(id) || document.querySelector(`[data-copilot-id="${id}"]`);
  if (!el) {
    throw new Error(`Element with id ${id} not found`);
  }
  (el as HTMLElement).click();
  return `Clicked element ${id}`;
}

export function fillInput(id: string, value: string): string {
  const el = document.getElementById(id) || document.querySelector(`[data-copilot-id="${id}"]`);
  if (!el) {
    throw new Error(`Input with id ${id} not found`);
  }
  
  const inputEl = el as HTMLInputElement;
  inputEl.value = value;
  
  // Dispatch events so React/Vue/Angular notice the change
  inputEl.dispatchEvent(new Event('input', { bubbles: true }));
  inputEl.dispatchEvent(new Event('change', { bubbles: true }));
  
  return `Filled input ${id} with value`;
}

export function scroll(direction: string): string {
  switch(direction) {
    case "down": window.scrollBy(0, window.innerHeight); break;
    case "up": window.scrollBy(0, -window.innerHeight); break;
    case "top": window.scrollTo(0, 0); break;
    case "bottom": window.scrollTo(0, document.body.scrollHeight); break;
    default: window.scrollBy(0, window.innerHeight);
  }
  return `Scrolled ${direction}`;
}

export function browserHistory(action: string): string {
  if (action === "back") {
    window.history.back();
    return "Navigating back";
  } else if (action === "forward") {
    window.history.forward();
    return "Navigating forward";
  } else if (action === "refresh") {
    window.location.reload();
    return "Refreshing page";
  }
  return "Unknown history action";
}

export function getCurrentUrl(): string {
  return window.location.href;
}
