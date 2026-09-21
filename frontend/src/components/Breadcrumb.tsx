// Breadcrumb component, aka the path of the page. ex. Documents > Github > TestProject2
import { ArrowRightToLineIcon } from 'lucide-react'

// Prop type structure
type Props = {
    currentPath: string
    // path:string -> void because the path is a string and the function doesn't return anything (void)
    onNavigate: (path: string) => void
}

// Main breadcrumb component
export default function Breadcrumb({ currentPath, onNavigate }: Props) {
    // Split segments
    const segments = currentPath.split('/').filter(Boolean)

    // Return component
    return (
        // Nav to create navigation section in main explorer
        <nav className="flex items-center gap-1 text-sm">
            {/* Home always navigates to base dir */}
            <span 
                onClick = {() => onNavigate('')}
                className="text-[#555] hover:text-[#e03030] cursor-pointer transition-colors duration-150"
            >
                Home
            </span>

            {/* Map segments to create breadcrumb links */}
            {segments.map((segment, index) => {
                // Build path to current segment, works by slicing segments to current index (+1 since exclusive) joining with '/'
                const pathToSegment = segments.slice(0, index + 1).join('/')
                const isLast = index === segments.length - 1

                // return span to navigate to each segment leading to current path
                return (
                    <span key = {pathToSegment} className='flex items-center gap-1'>
                        <ArrowRightToLineIcon className="text-[#333] w-4 h-4" />
                        <span 
                            onClick = {() => onNavigate(pathToSegment)} 
                            className={`cursor-pointer transition-colors duration-150 ${
                                isLast
                                    ? 'text-white'
                                    : 'text-[#555] hover:text-[#e03030]'
                            }`}
                        >
                            {segment}
                        </span>
                    </span>
                )
            })}
        </nav>
    )
}