const baseUrl = 'http://localhost:8080';

const chatWindow = document.getElementById('chat-window');
const recordButton = document.getElementById('record-btn');
recordButton.state = {
  initial: true,
  recording: false,
}

let mediaRecorder;
let audioChunks = [];

function setAuthorization(userId, userSecret) {
  localStorage.setItem('userId', userId);
  localStorage.setItem('userSecret', userSecret);
}

function getAuthorization() {
  const userId = localStorage.getItem('userId');
  const userSecret = localStorage.getItem('userSecret');

  // encode ke base64
  return btoa(`${userId}:${userSecret}`);
}

recordButton.onclick = () => {
  // pertama kali record button diklik
  if (recordButton.state.initial) {
    initChat();
    return;
  }

  // jika tidak sedang merekam
  if (!recordButton.state.recording) {
    startRecording();
    return;
  }

  // jika sedang merekam
  stopRecording();
}

async function initChat() {
  try {
    const response = await fetch(`${baseUrl}/chat/start`)
    const data = await response.json();

    // simpan userId dan userSecret
    const userId = data.data.id;
    const userSecret = data.data.secret;
    setAuthorization(userId, userSecret);

    // tampilkan pesan awal
    const initialMessage = data.data.text;
    appendMessage(initialMessage, 'assistant');

    // putar audio awal
    const initialAudio = data.data.audio;
    decodeAndPlayAudio(initialAudio);

    // atur button sudah diklik
    buttonIdle();
  } catch (error) {
    console.error("Error:", error);
    alert("Error starting chat, please try again.");
  }
}

async function startRecording() {
  try {
    // buat rekaman audio
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })

    // ubah button menjadi sedang merekam
    buttonRecording();

    mediaRecorder = new MediaRecorder(stream);
    mediaRecorder.start();

    // tambahkan audio dari mediaRecorder
    mediaRecorder.ondataavailable = (event) => {
      audioChunks.push(event.data);
    };

    // buat jadi blob
    mediaRecorder.onstop = () => {
      const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });

      // reset audio chunks
      audioChunks = [];

      // kirim audio ke server
      sendAudio(audioBlob);
    };
  } catch (error) {
    console.error("Error:", error);
    alert("Error starting recording, please try again.");

    // atur button menjadi menunggu merekam
    buttonIdle();
  }
}

async function sendAudio(audioBlob) {
  // buat form data
  const formData = new FormData();
  formData.append('file', audioBlob, 'audio.wav');

  try {
    // kirim audio ke server
    const response = await fetch(`${baseUrl}/chat/answer`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Basic ${getAuthorization()}`
      }
    })
    const data = await response.json()
    // tampilkan pesan hasil transkripsi
    const userMessage = data.data.prompt.text;
    appendMessage(userMessage, 'user');

    // tampikan pesan jawaban
    const replyMessage = data.data.answer.text;
    appendMessage(replyMessage, 'assistant')

    // putar audio jawaban
    const replyAudio = data.data.answer.audio;
    decodeAndPlayAudio(replyAudio);

    // atur button menjadi menunggu merekam
    buttonIdle();
  } catch (error) {
    console.error('Error:', error);
    alert('Error processing answer, please try again.');

    // atur button menjadi menunggu merekam
    buttonIdle();
  }

}

function stopRecording() {
  // hentikan rekaman
  mediaRecorder.stop();

  // ubah button menjadi sedang diproses
  buttonProcessing();
}

function appendMessage(message, type) {
  // buat div
  const messageDiv = document.createElement('div');

  // atur class dan isi
  messageDiv.className = `message ${type}`;
  messageDiv.textContent = message;

  // tambahkan ke chat window
  chatWindow.appendChild(messageDiv);

  // scroll ke bawah
  chatWindow.scrollTop = chatWindow.scrollHeight;
}

function decodeAndPlayAudio(encodedAudio) {
  // decode base64
  const audioData = atob(encodedAudio);

  // ubah menjadi array buffer
  const audioBuffer = new ArrayBuffer(audioData.length);
  const audioView = new Uint8Array(audioBuffer);

  // isi array buffer dengan data audio
  for (let i = 0; i < audioData.length; i++) {
    audioView[i] = audioData.charCodeAt(i);
  }

  // buat blob dari array buffer
  const audioBlob = new Blob([audioBuffer], { type: 'audio/wav' });

  // buat url dari blob
  const audioUrl = URL.createObjectURL(audioBlob);

  // putar audio
  const audio = new Audio(audioUrl);
  audio.play();
}

function buttonRecording() {
  // atur state button
  recordButton.state.recording = true;

  // atur text button
  recordButton.innerHTML = '<i class="bi bi-stop-fill"></i> Save Answer';
}

function buttonIdle() {
  // atur button agar bisa diklik
  recordButton.disabled = false;

  // atur state button
  recordButton.state.initial = false;
  recordButton.state.recording = false;

  // atur text button
  recordButton.innerHTML = '<i class="bi bi-record-circle"></i> Record Answer';
}

function buttonProcessing() {
  // atur button agar tidak bisa diklik & beri loading spinner
  recordButton.disabled = true;
  recordButton.innerHTML = '<span class="spinner-border spinner-border-sm" aria-hidden="true"></span><span role="status"> Processing...</span>'
}