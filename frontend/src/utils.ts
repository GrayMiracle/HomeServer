// Utilities for file helpers


// Function to convert bytes to readable format
export function formatSize(bytes: number): string {
    if (bytes >= 1099511627776) return (bytes / 1099511627776).toFixed(2) + ' TB'
    if (bytes >= 1073741824) return (bytes / 1073741824).toFixed(2) + ' GB'
    if (bytes >= 1048576) return (bytes / 1048576).toFixed(2) + ' MB'
    if (bytes >= 1024) return (bytes / 1024).toFixed(2) + ' KB'
    return bytes + ' B'
}

// Function duration to readable format
export function formatDuration(seconds: number): string {
    const hrs = Math.floor(seconds / 3600)
    const mins = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    // If there are hours, start with them then pad minutes and seconds
    if (hrs > 0) return `${hrs}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    // Otherwise start with mins and pad seconds
    return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Function to format date to readable format
export function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    })
}