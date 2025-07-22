from pyannote.audio import Pipeline
import torch
import os

class DiarizationService:
    def __init__(self, auth_token=None, device="cpu"):
        if not auth_token:
            print("⚠️ Hugging Face token not set. Diarization will not be available.")
            self.pipeline = None
            return

        print(f"   Loading Pyannote Diarization model on device '{device}'...")
        try:
            self.pipeline = Pipeline.from_pretrained(
                "pyannote/speaker-diarization-3.1",
                use_auth_token=auth_token
            )
            if device == "gpu" and torch.cuda.is_available():
                self.pipeline = self.pipeline.to(torch.device("cuda"))
            print("   Pyannote model loaded successfully.")
        except Exception as e:
            print(f"   Error loading Pyannote model: {e}")
            self.pipeline = None

    def diarize(self, audio_path: str):
        if not self.pipeline:
            raise Exception("Diarization pipeline is not available.")

        print("   Starting diarization...")
        diarization_result = self.pipeline(audio_path)
        print("   Diarization finished.")
        return diarization_result