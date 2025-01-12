#!/usr/bin/env bash

# Paths
OPERATOR_DIR=$(pwd)
GITOPS_DIR="$OPERATOR_DIR/../cointrac-gitops"
BASE_DIR="$GITOPS_DIR/base"

# Ensure the GitOps directory exists
if [ ! -d "$GITOPS_DIR" ]; then
	echo "Error: GitOps repository not found at $GITOPS_DIR"
	exit 1
fi

echo "Syncing files from $OPERATOR_DIR to $GITOPS_DIR..."

# Create necessary directories in GitOps repo
mkdir -p "$BASE_DIR/crds"
mkdir -p "$BASE_DIR/rbac"
mkdir -p "$BASE_DIR/operator"

# Copy CRDs
echo "Copying CRDs..."
CRD_FILES=()
for file in "$OPERATOR_DIR/config/crd/bases/"*.yaml; do
	cp -v "$file" "$BASE_DIR/crds/"
	CRD_FILES+=("crds/$(basename "$file")")
done

# Copy RBAC files
echo "Copying RBAC files..."
RBAC_FILES=()
for file in "$OPERATOR_DIR/config/rbac/"*.yaml; do
	cp -v "$file" "$BASE_DIR/rbac/"
	RBAC_FILES+=("rbac/$(basename "$file")")
done

# Copy Operator Deployment
echo "Copying Operator Deployment..."
cp -v "$OPERATOR_DIR/config/manager/manager.yaml" "$BASE_DIR/operator/operator-deployment.yaml"

# Generate kustomization.yaml dynamically
echo "Generating kustomization.yaml..."
cat >"$BASE_DIR/kustomization.yaml" <<EOL
resources:
EOL

# Add CRD files to kustomization.yaml
for file in "${CRD_FILES[@]}"; do
	echo "- $file" >>"$BASE_DIR/kustomization.yaml"
done

# Add RBAC files to kustomization.yaml
for file in "${RBAC_FILES[@]}"; do
	echo "- $file" >>"$BASE_DIR/kustomization.yaml"
done

# Add Operator Deployment to kustomization.yaml
echo "- operator/operator-deployment.yaml" >>"$BASE_DIR/kustomization.yaml"

echo "Sync completed successfully. All resources have been added to $GITOPS_DIR/base/kustomization.yaml."
