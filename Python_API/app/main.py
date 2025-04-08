from fastapi import FastAPI, UploadFile, HTTPException
from .predict import predict_image

app = FastAPI(title="CNN Prediction API")

@app.get("/test")
async def test_endpoint():
    """Simple test endpoint to verify communication."""
    return {"message": "If you are seeing this you are in the clear!"}

@app.post("/predict")
async def predict(file: UploadFile):
    """Receive an image and return the prediction vector."""
    try:
        # Read the image file
        image_data = await file.read()
        if not image_data:
            raise HTTPException(status_code=400, detail="No image data provided")
        
        # Get predictions
        predictions = predict_image(image_data)
        
        # Return the vector as JSON
        return {"predictions": predictions}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Prediction failed: {str(e)}")
