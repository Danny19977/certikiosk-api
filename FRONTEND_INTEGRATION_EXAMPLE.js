/**
 * Frontend Integration Example for Google Drive Downloads
 * 
 * This file shows how to update your DocumentViewer.js to use the backend proxy
 * endpoints instead of direct Google Drive URLs (which cause CORS errors).
 */

// ============================================
// CONFIGURATION
// ============================================

const API_BASE_URL = 'http://localhost:8081'; // Update this for production

// ============================================
// HELPER FUNCTIONS
// ============================================

/**
 * Extract Google Drive file ID from various URL formats
 * @param {string} url - Google Drive URL
 * @returns {string} - File ID
 */
function extractGoogleDriveFileId(url) {
  if (!url) return null;

  // Already a file ID (no slashes or protocol)
  if (!url.includes('/') && !url.includes('http')) {
    return url;
  }

  // https://drive.google.com/file/d/FILE_ID/view
  let match = url.match(/\/file\/d\/([a-zA-Z0-9_-]+)/);
  if (match) return match[1];

  // https://drive.google.com/open?id=FILE_ID
  match = url.match(/[?&]id=([a-zA-Z0-9_-]+)/);
  if (match) return match[1];

  // https://drive.google.com/uc?export=download&id=FILE_ID
  match = url.match(/[?&]id=([a-zA-Z0-9_-]+)/);
  if (match) return match[1];

  return null;
}

// ============================================
// MAIN DOWNLOAD FUNCTION (USE THIS)
// ============================================

/**
 * Download Google Drive document via backend proxy (NO CORS ISSUES!)
 * @param {string} fileIdOrUrl - Google Drive file ID or URL
 * @returns {Promise<Blob>} - PDF blob
 */
async function downloadGoogleDriveDocument(fileIdOrUrl) {
  try {
    console.log('🔄 Starting download via backend proxy...');
    
    // Extract file ID from URL if needed
    const fileId = extractGoogleDriveFileId(fileIdOrUrl);
    
    if (!fileId) {
      throw new Error('Invalid Google Drive file ID or URL');
    }

    console.log(`📄 File ID: ${fileId}`);

    // Use backend proxy endpoint (NO CORS!)
    const proxyUrl = `${API_BASE_URL}/api/public/documents/gdrive/download/${fileId}`;
    
    console.log(`🔗 Proxy URL: ${proxyUrl}`);

    const response = await fetch(proxyUrl, {
      method: 'GET',
      headers: {
        'Accept': 'application/pdf',
      }
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
    }

    const blob = await response.blob();
    
    console.log('✅ Download successful!');
    console.log(`📦 Blob size: ${blob.size} bytes`);
    console.log(`📋 Blob type: ${blob.type}`);

    return blob;

  } catch (error) {
    console.error('❌ Download failed:', error);
    throw error;
  }
}

// ============================================
// GET FILE METADATA
// ============================================

/**
 * Get Google Drive file metadata
 * @param {string} fileIdOrUrl - Google Drive file ID or URL
 * @returns {Promise<Object>} - File metadata
 */
async function getGoogleDriveFileMetadata(fileIdOrUrl) {
  try {
    const fileId = extractGoogleDriveFileId(fileIdOrUrl);
    
    if (!fileId) {
      throw new Error('Invalid Google Drive file ID or URL');
    }

    const metadataUrl = `${API_BASE_URL}/api/public/documents/gdrive/metadata/${fileId}`;
    
    const response = await fetch(metadataUrl);

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    
    if (result.status === 'success') {
      return result.data;
    } else {
      throw new Error(result.message || 'Failed to get metadata');
    }

  } catch (error) {
    console.error('❌ Metadata fetch failed:', error);
    throw error;
    }
}

// ============================================
// PRINT FUNCTION
// ============================================

/**
 * Print Google Drive PDF document
 * @param {string} fileIdOrUrl - Google Drive file ID or URL
 */
async function printGoogleDriveDocument(fileIdOrUrl) {
  try {
    console.log('🖨️ Preparing document for printing...');
    
    const blob = await downloadGoogleDriveDocument(fileIdOrUrl);
    
    // Create object URL for the blob
    const blobUrl = URL.createObjectURL(blob);
    
    // Open in new window and trigger print
    const printWindow = window.open(blobUrl, '_blank');
    
    if (printWindow) {
      printWindow.onload = function() {
        printWindow.print();
        
        // Clean up after printing
        setTimeout(() => {
          URL.revokeObjectURL(blobUrl);
        }, 100);
      };
    } else {
      console.warn('⚠️ Popup blocked. Please allow popups for this site.');
      // Fallback: download the file
      downloadFile(blob, 'document.pdf');
    }

    console.log('✅ Print dialog opened');

  } catch (error) {
    console.error('❌ Print failed:', error);
    alert('Failed to print document. Please try again.');
  }
}

// ============================================
// EMAIL FUNCTION
// ============================================

/**
 * Send Google Drive document via email
 * @param {string} fileIdOrUrl - Google Drive file ID or URL
 * @param {string} email - Recipient email
 * @param {string} documentType - Type of document (optional)
 * @param {string} documentName - Name of document (optional)
 */
async function emailGoogleDriveDocument(fileIdOrUrl, email, documentType = 'Document', documentName = 'document') {
  try {
    console.log('📧 Sending document via email...');
    
    const fileId = extractGoogleDriveFileId(fileIdOrUrl);
    
    if (!fileId) {
      throw new Error('Invalid Google Drive file ID or URL');
    }

    if (!email || !email.includes('@')) {
      throw new Error('Invalid email address');
    }

    const emailUrl = `${API_BASE_URL}/api/public/documents/send-email-gdrive`;
    
    const response = await fetch(emailUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email,
        file_id: fileId,
        document_type: documentType,
        document_name: documentName
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    
    if (result.status === 'success') {
      console.log('✅ Email sent successfully to:', email);
      return result;
    } else {
      throw new Error(result.message || 'Failed to send email');
    }

  } catch (error) {
    console.error('❌ Email sending failed:', error);
    throw error;
  }
}

// ============================================
// HELPER: DOWNLOAD FILE TO DISK
// ============================================

/**
 * Download blob as file
 * @param {Blob} blob - File blob
 * @param {string} filename - Filename
 */
function downloadFile(blob, filename) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

// ============================================
// USAGE EXAMPLES
// ============================================

/**
 * Example 1: Download and display PDF
 */
async function example1_DownloadAndDisplay() {
  const fileId = '1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_';
  
  try {
    const blob = await downloadGoogleDriveDocument(fileId);
    const blobUrl = URL.createObjectURL(blob);
    
    // Display in iframe
    document.getElementById('pdf-viewer').src = blobUrl;
    
  } catch (error) {
    console.error('Error:', error);
  }
}

/**
 * Example 2: Get file metadata
 */
async function example2_GetMetadata() {
  const fileId = '1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_';
  
  try {
    const metadata = await getGoogleDriveFileMetadata(fileId);
    console.log('File metadata:', metadata);
    // { id, name, mimeType, size, downloadUrl, viewUrl, proxy_url, ... }
    
  } catch (error) {
    console.error('Error:', error);
  }
}

/**
 * Example 3: Print document
 */
async function example3_Print() {
  const fileId = '1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_';
  
  try {
    await printGoogleDriveDocument(fileId);
  } catch (error) {
    console.error('Error:', error);
  }
}

/**
 * Example 4: Email document
 */
async function example4_Email() {
  const fileId = '1KHBHi5LBNmMlBVpti1WY18H6CJZasyN_';
  const email = 'user@example.com';
  
  try {
    await emailGoogleDriveDocument(fileId, email, 'Birth Certificate', 'birth_cert');
    alert('Email sent successfully!');
  } catch (error) {
    console.error('Error:', error);
    alert('Failed to send email');
  }
}

/**
 * Example 5: React Component Integration
 */
function DocumentViewerReactExample() {
  // Inside your React component:
  
  const handleDownload = async (googleDriveUrl) => {
    try {
      const blob = await downloadGoogleDriveDocument(googleDriveUrl);
      downloadFile(blob, 'document.pdf');
    } catch (error) {
      console.error('Download error:', error);
      alert('Failed to download document');
    }
  };

  const handlePrint = async (googleDriveUrl) => {
    try {
      await printGoogleDriveDocument(googleDriveUrl);
    } catch (error) {
      console.error('Print error:', error);
      alert('Failed to print document');
    }
  };

  const handleEmail = async (googleDriveUrl, email) => {
    try {
      await emailGoogleDriveDocument(googleDriveUrl, email);
      alert('Email sent successfully!');
    } catch (error) {
      console.error('Email error:', error);
      alert('Failed to send email');
    }
  };

  // Return your JSX with buttons that call these handlers
}

// ============================================
// EXPORT FOR MODULE SYSTEMS
// ============================================

// For ES6 modules
export {
  downloadGoogleDriveDocument,
  getGoogleDriveFileMetadata,
  printGoogleDriveDocument,
  emailGoogleDriveDocument,
  extractGoogleDriveFileId,
  downloadFile
};

// For CommonJS
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    downloadGoogleDriveDocument,
    getGoogleDriveFileMetadata,
    printGoogleDriveDocument,
    emailGoogleDriveDocument,
    extractGoogleDriveFileId,
    downloadFile
  };
}
