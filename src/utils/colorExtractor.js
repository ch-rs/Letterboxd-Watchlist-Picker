const colorThief = new ColorThief();

// Create a new file for color extraction
export function extractDominantColor(imgElement) {
    return new Promise((resolve) => {
        // Check if image is already loaded
        if (imgElement.complete && imgElement.naturalWidth !== 0) {
            processImage(imgElement);
        } else {
            imgElement.onload = () => processImage(imgElement);
        }

        function processImage(img) {
            const dominantColors = colorThief.getPalette(img, 10)
            let dominantColor = dominantColors[0]
            for (let scheme in dominantColors) {
                if (scheme[0] + scheme[1] + scheme[2] > dominantColor[0] + dominantColor[1] + dominantColor[2]) {
                    scheme = dominantColor
                }
            }

            resolve(dominantColor);
        }

        imgElement.onerror = () => {
            resolve([255, 255, 255]); // Default white on error
        };
    });
} 