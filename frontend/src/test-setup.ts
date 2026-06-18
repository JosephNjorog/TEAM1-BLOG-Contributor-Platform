// jsdom has no real layout engine, so Range/Element don't implement
// getBoundingClientRect - stub it out so code that calls it (purely for
// positioning a UI popover, never for the offset math itself) doesn't
// throw in tests.
if (typeof Range !== "undefined") {
  Range.prototype.getBoundingClientRect = function () {
    return { x: 0, y: 0, width: 0, height: 0, top: 0, right: 0, bottom: 0, left: 0, toJSON: () => ({}) } as DOMRect;
  };
}
