// Computes a selected range as plain-text character offsets relative to a
// container element, by walking its text nodes. This matches how the
// backend's word counter treats content (HTML stripped, plain text), so a
// suggestion's rangeStart/rangeEnd stay meaningful without needing to track
// raw HTML offsets across re-renders.
export interface SelectionRange {
  start: number;
  end: number;
  text: string;
  rect: DOMRect;
}

export function captureSelection(container: HTMLElement): SelectionRange | null {
  const selection = window.getSelection();
  if (!selection || selection.isCollapsed || selection.rangeCount === 0) return null;

  const range = selection.getRangeAt(0);
  if (!container.contains(range.commonAncestorContainer)) return null;

  const text = selection.toString();
  if (text.trim() === "") return null;

  const offsetOf = (node: Node, nodeOffset: number): number => {
    const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT);
    let total = 0;
    let current = walker.nextNode();
    while (current) {
      if (current === node) return total + nodeOffset;
      total += current.textContent?.length ?? 0;
      current = walker.nextNode();
    }
    return total;
  };

  const start = offsetOf(range.startContainer, range.startOffset);
  const end = offsetOf(range.endContainer, range.endOffset);
  const rect = range.getBoundingClientRect();

  return { start, end, text, rect };
}
