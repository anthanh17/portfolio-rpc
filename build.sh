ID="338511574097"
REGION="ap-southeast-1"
REPOSITORY_NAME=`git remote show origin | grep 'Fetch URL' | sed 's/.*\/\([^\/]*\)\.git/\1/'`
PREFIX_AWS="${ID}.dkr.ecr.${REGION}.amazonaws.com"

# Login Repo
aws ecr get-login-password --region ap-southeast-1 --profile rd | docker login --username AWS --password-stdin 338511574097.dkr.ecr.ap-southeast-1.amazonaws.com

# Check if the ECR repository exists
if aws ecr describe-repositories --repository-names "${REPOSITORY_NAME}" --region "${REGION}" --profile rd > /dev/null 2>&1; then
  echo "Repos>\tECR repository '${REPOSITORY_NAME}' already exists"
else
  # Create the ECR repository
  aws ecr create-repository --repository-name "${REPOSITORY_NAME}" --region "${REGION}" --profile rd
  echo "Repos>\tECR repository '${REPOSITORY_NAME}' created"
fi

# Get current branch
CUR_BRANCH=`git rev-parse --abbrev-ref HEAD | sed 's/\//\-/g'`
# Get Commit ID
COMMIT_ID=`git log | head -1 | sed s/'commit '//`
SUB_COMMIT_ID=$(echo $COMMIT_ID | cut -c 1-7)
# Get Commit date
COMMIT_DATE_FORMAT="$(date -d "@$(git --git-dir ".git" show -s --format=%ct HEAD)" +"%y.%m%d.%H%M")"
# Image name
IMAGE="${PREFIX_AWS}/${REPOSITORY_NAME}:${COMMIT_DATE_FORMAT}_${CUR_BRANCH}_${SUB_COMMIT_ID}"

echo "Building image $IMAGE\n"

# Build Image
docker build . -t $IMAGE

# Push Image
docker push $IMAGE
