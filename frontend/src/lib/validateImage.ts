const MAX_BYTES = 5 * 1024 * 1024;
const MIN_SIDE = 1360;

// Mirrors internal/cloudinary/validate.go so a bad file gets caught before
// the upload round-trip instead of after.
export function validateBannerFile(file: File): Promise<string | null> {
  if (!["image/jpeg", "image/png"].includes(file.type)) {
    return Promise.resolve("Only JPG and PNG images are accepted.");
  }
  if (file.size > MAX_BYTES) {
    return Promise.resolve("File exceeds the 5MB maximum.");
  }
  return new Promise((resolve) => {
    const img = new Image();
    const url = URL.createObjectURL(file);
    img.onload = () => {
      URL.revokeObjectURL(url);
      if (img.naturalWidth < MIN_SIDE || img.naturalHeight < MIN_SIDE) {
        resolve(`Image must be at least ${MIN_SIDE}x${MIN_SIDE}px (this one is ${img.naturalWidth}x${img.naturalHeight}).`);
      } else {
        resolve(null);
      }
    };
    img.onerror = () => {
      URL.revokeObjectURL(url);
      resolve("Could not read this image file.");
    };
    img.src = url;
  });
}
