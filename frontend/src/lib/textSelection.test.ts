import { afterEach, describe, expect, it } from "vitest";
import { captureSelection } from "./textSelection";

function selectText(container: HTMLElement, startNode: Node, startOffset: number, endNode: Node, endOffset: number) {
  const range = document.createRange();
  range.setStart(startNode, startOffset);
  range.setEnd(endNode, endOffset);
  const selection = window.getSelection();
  selection?.removeAllRanges();
  selection?.addRange(range);
  return range;
}

describe("captureSelection", () => {
  afterEach(() => {
    window.getSelection()?.removeAllRanges();
    document.body.innerHTML = "";
  });

  it("computes plain-text offsets within a single text node", () => {
    const container = document.createElement("div");
    container.innerHTML = "<p>The quick brown fox jumps</p>";
    document.body.appendChild(container);

    const textNode = container.querySelector("p")!.firstChild!;
    selectText(container, textNode, 4, textNode, 15); // "quick brown"

    const result = captureSelection(container);
    expect(result).not.toBeNull();
    expect(result?.start).toBe(4);
    expect(result?.end).toBe(15);
    expect(result?.text).toBe("quick brown");
  });

  it("sums offsets correctly across multiple text nodes (inline markup)", () => {
    const container = document.createElement("div");
    container.innerHTML = "<p>Hello <strong>world</strong> there</p>";
    document.body.appendChild(container);

    const p = container.querySelector("p")!;
    const worldNode = p.querySelector("strong")!.firstChild!; // "world"
    const thereNode = p.lastChild!; // " there"
    selectText(container, worldNode, 0, thereNode, 6); // "world there"

    const result = captureSelection(container);
    expect(result).not.toBeNull();
    expect(result?.start).toBe(6); // length of "Hello "
    expect(result?.end).toBe(17); // length of "Hello " + "world" + " there"
    expect(result?.text).toBe("world there");
  });

  it("returns null when there is no selection", () => {
    const container = document.createElement("div");
    container.innerHTML = "<p>Nothing selected here</p>";
    document.body.appendChild(container);

    expect(captureSelection(container)).toBeNull();
  });

  it("returns null when the selection is collapsed (just a cursor)", () => {
    const container = document.createElement("div");
    container.innerHTML = "<p>Some text</p>";
    document.body.appendChild(container);

    const textNode = container.querySelector("p")!.firstChild!;
    selectText(container, textNode, 3, textNode, 3);

    expect(captureSelection(container)).toBeNull();
  });

  it("returns null when the selection is outside the container", () => {
    const container = document.createElement("div");
    container.innerHTML = "<p>Inside</p>";
    const outside = document.createElement("p");
    outside.textContent = "Outside the container";
    document.body.append(container, outside);

    selectText(container, outside.firstChild!, 0, outside.firstChild!, 7);

    expect(captureSelection(container)).toBeNull();
  });
});
