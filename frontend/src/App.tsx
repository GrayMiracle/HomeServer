import './App.css'
import { Route, Routes, Navigate } from 'react-router-dom'
import FileExplorer from './pages/FileExplorer'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import AuthGate from './components/AuthGate'
import GuestGate from './components/GuestGate'

function App() {
	return (
		<div className = 'min-h-screen bg-black text-white'>
			<Routes>
				<Route path = '/login' element = {<GuestGate><Login /></GuestGate>} />
				<Route path = '/dashboard' element = {<AuthGate><Dashboard /></AuthGate>} />
				<Route path = '/' element = {<AuthGate><FileExplorer /></AuthGate>} />
				<Route path = '*' element = {<Navigate to = '/' replace />} />
			</Routes>
		</div>
	)
}
export default App
