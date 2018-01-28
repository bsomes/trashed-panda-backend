mkdir lib
TF_TYPE="cpu" # Change to "gpu" for GPU support
 TARGET_DIRECTORY='./lib'
 curl -L \
   "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-linux-x86_64-1.5.0.tar.gz" |
tar -C $TARGET_DIRECTORY -xz

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname $CURRENT_DIR)"

export LIBRARY_PATH=$LIBRARY_PATH:$PARENT_DIR/lib/lib
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PARENT_DIR/lib/lib