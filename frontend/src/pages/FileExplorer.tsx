import { useState, useEffect } from 'react'
import type { FileItem } from '../types'
import { API } from '../config'
import { Upload, X } from 'lucide-react'
import { BarChart2 } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { apiFetch } from '../api'

// Component imports
import VideoPlayer from '../components/VideoPlayer'
import AudioPlayer from '../components/AudioPlayer'
import ImageViewer from '../components/ImageViewer'
import UploadModal from '../components/UploadModal'
import PinModal from '../components/PinModal'
import FileCard from '../components/FileCard'
import Breadcrumb from '../components/Breadcrumb'
import TextViewer from '../components/TextViewer'
import SettingsMenu from '../components/Settings'

// Main file explorer
export default function FileExplorer() {
    const navigate = useNavigate()
    // useStates
    const [loading, setLoading] = useState(false)
    const [files, setFiles] = useState<FileItem[]>([])
    const [currentPath, setCurrentPath] = useState('')
    const [error, setError] = useState('')

    // selected media
    const [selectedVideo, setSelectedVideo] = useState<FileItem | null>(null)
    const [selectedAudio, setSelectedAudio] = useState<FileItem | null>(null)
    const [selectedImage, setSelectedImage] = useState<FileItem | null>(null)
    const [selectedPdf, setSelectedPdf] = useState<FileItem | null>(null)
    const [selectedText, setSelectedText] = useState<FileItem | null>(null)
    const [selectedDownload, setSelectedDownload] = useState<FileItem | null>(null)
    const [showUpload, setShowUpload] = useState(false)
    const [pinAction, setPinAction] = useState<(() => void) | null>(null)

    // useEffect to get new files on path chagne
    useEffect(() => {
        fetchFiles()
    }, [currentPath])

    // actual function to fetch files
    async function fetchFiles() {
        console.log('Fetching path:', currentPath)
        setLoading(true)
        setError('')
        try {
            const res = await apiFetch(`/files?path=${encodeURIComponent(currentPath)}`)
            if (!res.ok) throw new Error('Failed to fetch')
            const data = await res.json()
            setFiles(Array.isArray(data) ? data : [])
        } catch (err) {
            setError('Failed to reach route')
        } finally {
            setLoading(false)
        }
    }

    // function to change path for folders
    function folderClick(folder: FileItem) {
        setCurrentPath(folder.path)
    }

    // function for pin required actions
    function requirePin(action: () => void) {
        setPinAction(() => action)
    }

    // function to delete files
    async function deleteFile(file: FileItem) {
        if (!window.confirm(`Are you sure you want to delete ${file.name}?`)) return
        requirePin(async () => {
            setError('')
            try {
                await apiFetch(`/file/delete?path=${encodeURIComponent(file.path)}`, { method: 'DELETE' })
                fetchFiles()
            } catch {
                setError('Could not delete')
            }
        })
    }

    // return content
    return (
        <div className = 'min-h-screen bg-black'>

            {/* Top bar */}
            <div className = 'flex items-center justify-between px-6 py-3 bg-[#111] border-b border-[#1f1f1f]'>
                {/* Breadcrumb nav */}
                <Breadcrumb currentPath = {currentPath} onNavigate = {setCurrentPath}/>

                {/* Right-side controls */}
                <div className="flex items-center gap-3">

                    {/* Dashboard */}
                    <button
                        onClick={() => navigate('/dashboard')}
                        className="flex items-center gap-2 bg-[#1f1f1f] hover:bg-[#2a2a2a] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
                    >
                        <BarChart2 size={14} />
                    </button>

                    {/* Upload */}
                    <button
                        onClick={() => requirePin(() => setShowUpload(true))}
                        className="flex items-center gap-2 bg-[#e03030] hover:bg-[#ff3c3c] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
                    >
                        <Upload size={14} /> Upload
                    </button>

                    {/* Settings */}
                    <SettingsMenu />
                </div>

                {showUpload && (
                    <UploadModal
                        currentPath = {currentPath}
                        onClose = {() => setShowUpload(false)}
                        onUploadSuccess = {fetchFiles}
                    />
                )}

            </div>
            

            {/* Statuses */}
            {loading && <p className="px-6 py-3 text-sm text-[#666]">Loading...</p>}
            {error && <p className="px-6 py-3 text-sm text-[#e03030]">{error}</p>}
            {pinAction && (
                <PinModal
                    onSuccess={() => {
                        const action = pinAction
                        setPinAction(null)
                        action()
                    }}
                    onCancel={() => setPinAction(null)}
                />
            )}

            {/* Array traversing through all files in current dir to render a filecard for each */}
            <div className = 'grid grid-cols-[repeat(auto-fill,minmax(160px,1fr))] gap-3 p-6'>
                {files.map(file => (
                    <FileCard
                        key = {file.path}
                        file = {file}
                        onFolderClick = {folderClick}
                        onVideoClick = {setSelectedVideo}
                        onImageClick = {setSelectedImage}
                        onAudioClick = {setSelectedAudio}
                        onPdfClick = {setSelectedPdf}
                        onTextClick = {setSelectedText}
                        onDownloadClick = {setSelectedDownload}
                        onDelete = {deleteFile}
                    />
                ))}
            </div>

            {/* Selected media display */}

            {/* Video: */}
            {selectedVideo && (
                <VideoPlayer file = {selectedVideo} onClose = {() => setSelectedVideo(null)}/>
            )}

            {/* Audio: */}
            {selectedAudio && (
                <AudioPlayer file = {selectedAudio} onClose = {() => setSelectedAudio(null)}/>
            )}

            {/* Image: */}
            {selectedImage && (
                <ImageViewer file = {selectedImage} onClose = {() => setSelectedImage(null)}/>
            )}

            {/* Pdf */}
            {selectedPdf && (
                <div className="fixed inset-0 bg-[#0a0a0a] z-50 flex flex-col">
                    {/* Top bar */}
                    <div className="flex items-center justify-between px-6 py-3 bg-[#111] border-b border-[#1f1f1f] shrink-0">
                        <Breadcrumb
                            currentPath={selectedPdf.path}
                            onNavigate={(path) => {
                                setCurrentPath(path)
                                setSelectedPdf(null)
                            }}
                        />
                        <button
                            onClick={() => setSelectedPdf(null)}
                            className="flex items-center gap-2 bg-[#1f1f1f] hover:bg-[#2a2a2a] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
                        >
                            <X size={14} />
                        </button>
                    </div>
                    {/* PDF iframe */}
                    <iframe
                        src={`${API}/file/download?path=${encodeURIComponent(selectedPdf.path)}&intent=view`}
                        className="flex-1 w-full border-none"
                    />
                </div>
            )}

            {/* Text */}
            {selectedText && (
                <TextViewer
                    file={selectedText}
                    onClose={() => setSelectedText(null)}
                    onNavigate={setCurrentPath}
                />
            )}

            {/* Other */}
            {selectedDownload && (
                <div>
                    <button onClick = {() => setSelectedDownload(null)}>Close</button>
                    {selectedDownload.name} {selectedDownload.size}
                    <a
                        href = {`${API}/file/download?path=${encodeURIComponent(selectedDownload.path)}&intent=download`}
                        download = {selectedDownload.name}
                        ref = {el => {
                            if (el) {
                                el.click()
                                setSelectedDownload(null)
                            }
                        }}
                    />
                </div>
            )}

        </div>
    )
}