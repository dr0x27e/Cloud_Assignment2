import tensorflow as tf
model = tf.keras.models.load_model("/home/drt/tempCloud/API_service/model/custom_model.keras", compile=False)
print(model.get_config())
