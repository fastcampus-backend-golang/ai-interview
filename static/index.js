const chatWindow = document.getElementById('chat-window');
const startBtn = document.getElementById('start-btn');
const recordBtn = document.getElementById('record-btn');

let mediaRecorder;
let audioChunks = [];

let userId = "";
let userSecret = ""

startBtn.addEventListener('click', () => {
    fetch('http://localhost:8080/chat/start')
        .then(response => response.json())
        .then(data => {
            userId = data.data.id;
            userSecret = data.data.secret;

            const systemMessage = data.data.text;
            appendMessage(systemMessage, 'reply');

            const replyAudio = data.data.audio;
            decodeBase64Audio(replyAudio);
        })
        .catch(error => {
            console.error('Error:', error);
            window.alert('Error starting chat, please try again.');
        });

    startBtn.style.display = 'none';
    recordBtn.style.display = 'block';
});

recordBtn.addEventListener('click', () => {
    if (recordBtn.textContent === 'Record Audio') {
        startRecording();
    } else {
        stopRecording();
    }
});

function startRecording() {
    navigator.mediaDevices.getUserMedia({ audio: true })
        .then(stream => {
            mediaRecorder = new MediaRecorder(stream);
            mediaRecorder.start();
            enableStopButton();

            mediaRecorder.ondataavailable = event => {
                audioChunks.push(event.data);
            };

            mediaRecorder.onstop = () => {
                const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
                audioChunks = [];
                handleAudioBlob(audioBlob);
            };
        });
}

function stopRecording() {
    mediaRecorder.stop();
    disableRecordButton();
}

function handleAudioBlob(audioBlob) {
    const formData = new FormData();
    formData.append('file', audioBlob, 'audio.wav');

    // base64 encode id and secret
    const auth = btoa(`${userId}:${userSecret}`);

    // Send audio to speech-to-text API
    fetch('http://localhost:8080/chat/answer', {
        method: 'POST',
        body: formData,
        headers: {
            'Authorization': `Basic ${auth}`
        }
    })
    .then(response => response.json())
    .then(data => {
        const userMessage = data.data.prompt.text;
        appendMessage(userMessage, 'user');
        
        const replyMessage = data.data.answer.text;
        appendMessage(replyMessage, 'reply')

        const replyAudio = data.data.answer.audio;
        decodeBase64Audio(replyAudio);
        
        enableRecordButton();
    })
    .catch(error => {
        console.error('Error:', error);
        window.alert('Error processing answer, please try again.');

        enableRecordButton();
    });
}

function appendMessage(message, type) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = message;
    chatWindow.appendChild(messageDiv);
    chatWindow.scrollTop = chatWindow.scrollHeight;
}

function decodeBase64Audio(encodedAudio) {
    const audioData = atob(encodedAudio);
    const audioBuffer = new ArrayBuffer(audioData.length);
    const audioView = new Uint8Array(audioBuffer);
    for (let i = 0; i < audioData.length; i++) {
        audioView[i] = audioData.charCodeAt(i);
    }
    const audioBlob = new Blob([audioBuffer], { type: 'audio/wav' });
    const audioUrl = URL.createObjectURL(audioBlob);

    const audio = new Audio(audioUrl);
    audio.play();
}

function disableRecordButton() {
    recordBtn.disabled = true;
    recordBtn.textContent = 'Processing';
}

function enableRecordButton() {
    recordBtn.disabled = false;
    recordBtn.textContent = 'Record Audio';
}

function enableStopButton() {
    recordBtn.textContent = 'Stop Recording';
}
