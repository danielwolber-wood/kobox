// ==UserScript==
// @name         Page HTML Uploader
// @namespace    http://tampermonkey.net/
// @version      1.0
// @description  Upload current page title and HTML to localhost API
// @author       DLWW
// @match        *://*/*
// @grant        none
// ==/UserScript==

(function() {
    'use strict';

    const uploadButton = document.createElement('button');
    uploadButton.innerHTML = '📤 Upload Page';
    uploadButton.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 10000;
        background: #007AFF;
        color: white;
        border: none;
        padding: 12px 16px;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        cursor: pointer;
        box-shadow: 0 2px 10px rgba(0,0,0,0.2);
        transition: all 0.2s ease;
        display: none;
    `;

    // Show/hide button based on Control key state
    let isControlPressed = false;

    function updateButtonVisibility() {
        uploadButton.style.display = isControlPressed ? 'block' : 'none';
    }

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Meta' && !isControlPressed) {
            isControlPressed = true;
            updateButtonVisibility();
        }
    });

    document.addEventListener('keyup', (event) => {
        if (event.key === 'Meta' && isControlPressed) {
            isControlPressed = false;
            updateButtonVisibility();
        }
    });

    // Also handle when the window loses focus (to reset state)
    window.addEventListener('blur', () => {
        isControlPressed = false;
        updateButtonVisibility();
    });

    uploadButton.addEventListener('mouseenter', () => {
        uploadButton.style.background = '#0056D6';
        uploadButton.style.transform = 'translateY(-1px)';
    });

    uploadButton.addEventListener('mouseleave', () => {
        uploadButton.style.background = '#007AFF';
        uploadButton.style.transform = 'translateY(0)';
    });

    async function uploadPageData() {
        try {
            // Change button state to loading
            uploadButton.innerHTML = '⏳ Uploading...';
            uploadButton.disabled = true;

            // Prepare the data to match Go struct format
            const pageData = {
                html: document.documentElement.outerHTML,
                title: document.title
            };

            console.log('Uploading page data:', {
                title: pageData.title,
                htmlLength: pageData.html.length,
                url: window.location.href
            });

            const response = await fetch('https://localhost:12332/v2/api/upload/html', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(pageData)
            });

            console.log('Response status:', response.status);
            console.log('Response headers:', [...response.headers.entries()]);

            if (response.ok) {
                const responseText = await response.text();
                console.log('Success response:', responseText);

                // Success state
                uploadButton.innerHTML = '✅ Uploaded!';
                uploadButton.style.background = '#28a745';

                // Reset after 2 seconds
                setTimeout(() => {
                    uploadButton.innerHTML = '📤 Upload Page';
                    uploadButton.style.background = '#007AFF';
                    uploadButton.disabled = false;
                }, 2000);
            } else {
                const errorText = await response.text();
                console.error('Error response body:', errorText);
                throw new Error(`HTTP ${response.status}: ${response.statusText} - ${errorText}`);
            }

        } catch (error) {
            console.error('Upload failed:', error);

            uploadButton.innerHTML = '❌ Failed';
            uploadButton.style.background = '#dc3545';

            const errorDiv = document.createElement('div');
            errorDiv.style.cssText = `
                position: fixed;
                top: 70px;
                right: 20px;
                z-index: 10001;
                background: #dc3545;
                color: white;
                padding: 8px 12px;
                border-radius: 4px;
                font-size: 12px;
                max-width: 300px;
                word-wrap: break-word;
            `;
            errorDiv.textContent = error.message;
            document.body.appendChild(errorDiv);

            setTimeout(() => {
                if (errorDiv.parentNode) {
                    errorDiv.parentNode.removeChild(errorDiv);
                }
                uploadButton.innerHTML = '📤 Upload Page';
                uploadButton.style.background = '#007AFF';
                uploadButton.disabled = false;
            }, 5000);
        }
    }

    uploadButton.addEventListener('click', uploadPageData);

    document.body.appendChild(uploadButton);

})();