import whisper
import torch

class TranscriptionService:
    def __init__(self, model_name="base", device="cpu"):
        print(f"   Loading Whisper model '{model_name}' on device '{device}'...")
        try:
            self.model = whisper.load_model(model_name, device=device)
            print("   Whisper model loaded successfully.")
        except Exception as e:
            print(f"   Error loading Whisper model: {e}")
            self.model = None

    def transcribe(self, audio_path: str):
        if not self.model:
            raise Exception("Whisper model is not available.")

        use_fp16 = self.model.device.type == 'cuda'

        print("   Starting transcription...")
        result = self.model.transcribe(
            audio_path,
            word_timestamps=True,
            language='ja',
            fp16=use_fp16
        )
        print("   Transcription finished.")
        return result