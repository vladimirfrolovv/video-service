import React, { useEffect, useState } from 'react'

const BACKEND_URL = 'http://localhost:8080'

function App() {
    const [file, setFile] = useState(null)
    const [videos, setVideos] = useState([])

    // Получаем список файлов
    const fetchVideos = async () => {
        try {
            const res = await fetch(`${BACKEND_URL}/list`)
            const data = await res.json()
            setVideos(data)
        } catch (err) {
            console.error('Ошибка запроса /list:', err)
        }
    }

    // При первой загрузке компонента — запрашиваем список
    useEffect(() => {
        fetchVideos()
    }, [])

    // Обработчик загрузки видео
    const handleUpload = async (e) => {
        e.preventDefault()
        if (!file) return
        const formData = new FormData()
        formData.append('file', file)

        try {
            const res = await fetch(`${BACKEND_URL}/upload`, {
                method: 'POST',
                body: formData,
            })
            const text = await res.text()
            console.log('Ответ /upload:', text)
            alert(text)
            setFile(null)
            // Обновим список после загрузки
            fetchVideos()
        } catch (err) {
            console.error('Ошибка при загрузке файла:', err)
        }
    }

    return (
        <div style={{ maxWidth: 600, margin: '0 auto', fontFamily: 'Arial, sans-serif' }}>
            <h1>My Videos</h1>

            <form onSubmit={handleUpload}>
                <input
                    type="file"
                    accept="video/*"
                    onChange={(e) => setFile(e.target.files[0])}
                />
                <button type="submit">Загрузить</button>
            </form>

            <button onClick={fetchVideos}>Обновить список</button>

            <hr />

            {videos.map((filename) => (
                <div key={filename} style={{ marginBottom: 20 }}>
                    <h3>{filename}</h3>
                    <video
                        width="100%"
                        height="300"
                        controls
                        src={`${BACKEND_URL}/video/${filename}`}
                    >
                        Ваш браузер не поддерживает видео-тег.
                    </video>
                </div>
            ))}
        </div>
    )
}

export default App
