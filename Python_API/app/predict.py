import tensorflow as tf
import cv2
import numpy as np
import os

# Custom loss function (no decorator needed for loading)
def custom_loss(y_true, y_pred):
    # Splitting the tensors into components
    true_conf    = y_true[..., 0]  # (batch, 7, 7)               (Confidence)
    true_boxes   = y_true[..., 1:5]  # (batch, 7, 7, 4)          (Bounding boxes)
    true_classes = y_true[..., 5:]  # (batch, 7, 7, num_classes) (class list)
    
    pred_conf    = y_pred[..., 0]  # (batch, 7, 7)               (Same deal)         
    pred_boxes   = y_pred[..., 1:5]  # (batch, 7, 7, 4)          ----||----
    pred_classes = y_pred[..., 5:]  # (batch, 7, 7, num_classes) ----||----
    
    # Masks for cells with and without objects
    obj_mask   = tf.cast(true_conf > 0, tf.float32)
    noobj_mask = 1.0 - obj_mask
    
    # Parameters for weighting
    weight_noobj = 0.3   # Weight for confidence when no object
    weight_coord = 4.0   # Weight for bounding box coordinates
    
    # 1. Confidence loss (using MSE)
    conf_diff       = tf.square(true_conf - pred_conf)
    conf_loss_obj   = tf.reduce_sum(obj_mask * conf_diff)
    conf_loss_noobj = tf.reduce_sum(noobj_mask * conf_diff)
    conf_loss       = conf_loss_obj + weight_noobj * conf_loss_noobj
    
    # 2. Bounding box loss (using MSE, only for cells with objects)
    box_diff = tf.square(true_boxes - pred_boxes)
    box_loss = tf.reduce_sum(obj_mask[..., tf.newaxis] * box_diff)
    box_loss = weight_coord * box_loss
    
    # 3. Class loss (using cross-entropy, only for cells with objects)
    cross_entropy = tf.keras.losses.categorical_crossentropy(
        true_classes, pred_classes, from_logits=False
    )
    class_loss = tf.reduce_sum(obj_mask * cross_entropy)
    
    # Total loss
    total_loss = conf_loss + box_loss + class_loss
    return total_loss

# Load the model with custom objects
MODEL_PATH = os.path.join(os.path.dirname(__file__), "..", "model", "custom_model.keras")
if not os.path.exists(MODEL_PATH):
    raise FileNotFoundError(f"Model file not found at: {MODEL_PATH}")
model = tf.keras.models.load_model(MODEL_PATH, custom_objects={"custom_loss": custom_loss})

# Constants
GRID_SIZE = 7
NUM_CLASSES = 2  # Dog (Class1), Cat (Class2)
INPUT_SIZE = 224  # Model input size

def preprocess_image(image_data):
    """Preprocess the uploaded image for the model."""
    nparr = np.frombuffer(image_data, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    if img is None:
        raise ValueError("Failed to decode image")
    
    img_resized = cv2.resize(img, (INPUT_SIZE, INPUT_SIZE))
    img_normalized = img_resized / 255.0
    img_input = np.expand_dims(img_normalized, axis=0)
    return img_input

def process_predictions(predictions):
    """Convert model output to a filtered list with grid positions and class labels."""
    predictions = predictions[0]  # Remove batch dimension
    
    result_vector = []
    for i in range(GRID_SIZE):
        for j in range(GRID_SIZE):
            pred = predictions[i, j]
            probability = float(pred[0])
            
            # Filter: only include predictions with P > 0.65
            if probability > 0.65:
                class_dog = float(pred[5])
                class_cat = float(pred[6])
                class_label = "Dog" if class_dog > class_cat else "Cat"
                
                vector = {
                    "grid": [i, j],  # Add grid position
                    "probability": probability,
                    "xmin": float(pred[1]),
                    "ymin": float(pred[2]),
                    "xmax": float(pred[3]),
                    "ymax": float(pred[4]),
                    "class": class_label  # Simplified class label
                }
                result_vector.append(vector)
    
    return result_vector

def predict_image(image_data):
    """Process the image and return the prediction vector."""
    img_input = preprocess_image(image_data)
    predictions = model.predict(img_input)
    return process_predictions(predictions)
