import os
import uvicorn
import yaml
from fastapi import FastAPI, UploadFile, File

from services.transcription import TranscriptionService
from services.diarization import DiarizationService

# --- 1. 設定ファイルの読み込み ---
CONFIG_PATH = "config.yml"
config = {}
if os.path.exists(CONFIG_PATH):
    with open(CONFIG_PATH, 'r') as f:
        config = yaml.safe_load(f)

# --- 2. 初期設定 ---
app = FastAPI()
TEMP_AUDIO_PATH = "temp_audio.wav"

print("✅ Python Voice Analysis Service is starting...")

# --- 3. AIサービスの初期化 ---
whisper_model_size = config.get('models', {}).get('whisper_model_size', 'base')
hf_token = config.get('models', {}).get('hugging_face_token', None)
device = config.get('compute', {}).get('device', 'cpu')

transcription_service = TranscriptionService(model_name=whisper_model_size, device=device)
diarization_service = DiarizationService(auth_token=hf_token, device=device)

# --- 4. APIエンドポイント ---
@app.get("/")
def read_root():
    return {"message": "Voice Analysis Service is running."}

@app.post("/v1/analyze")
async def analyze_audio(audio_file: UploadFile = File(...)):
    if not transcription_service.model or not diarization_service.pipeline:
        return {"error": "AI models are not available."}

    print("   Receiving audio file for analysis...")

    with open(TEMP_AUDIO_PATH, "wb") as buffer:
        buffer.write(await audio_file.read())

    try:
        whisper_result = transcription_service.transcribe(TEMP_AUDIO_PATH)

        word_segments = []
        if 'segments' in whisper_result and whisper_result['segments']:
            for segment in whisper_result['segments']:
                if 'words' in segment:
                    word_segments.extend(segment['words'])

        if not word_segments:
            return {"error": "Whisper could not detect any words with timestamps."}

        diarization_result = diarization_service.diarize(TEMP_AUDIO_PATH)

        final_result = []
        current_speaker = ""
        current_text = ""

        for word_info in word_segments:
            word_start_time = word_info['start']

            speaker = "UNKNOWN"
            for turn, _, speaker_label in diarization_result.itertracks(yield_label=True):
                if turn.start <= word_start_time <= turn.end:
                    speaker = speaker_label
                    break

            if current_speaker != speaker and current_text:
                final_result.append({"speaker": current_speaker, "text": current_text.strip()})
                current_text = ""

            current_speaker = speaker
            current_text += word_info['word']

        if current_text:
            final_result.append({"speaker": current_speaker, "text": current_text.strip()})

        print(f"   Analysis result: {final_result}")

    except Exception as e:
        print(f"   Error during processing: {e}")
        return {"error": f"Processing failed: {e}"}
    finally:
        if os.path.exists(TEMP_AUDIO_PATH):
            os.remove(TEMP_AUDIO_PATH)

    return {"segments": final_result}

# --- 5. サーバーの起動 ---
if __name__ == "__main__":
    port = config.get('server_ports', {}).get('voice_analysis_service', 8000)
    uvicorn.run(app, host="0.0.0.0", port=port)