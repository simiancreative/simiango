currversion=$(git describe --abbrev=0 --tags)
echo "current -> $currversion"

version=$(./bin/increment-version.sh -p $currversion)
echo "new -> $version"

git tag $version
git push origin $version
