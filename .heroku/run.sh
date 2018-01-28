mkdir 'lib'
TF_TYPE="cpu" # Change to "gpu" for GPU support
 TARGET_DIRECTORY="/usr/local"
 curl -L \
   "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-linux-x86_64-1.5.0.tar.gz" |
tar -C $TARGET_DIRECTORY -xz

#export LIBRARY_PATH=$LIBRARY_PATH:./lib
#export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./lib