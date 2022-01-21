currversion=$(git describe --abbrev=0 --tags)
echo "current -> $currversion"

newVersion=$(./bin/increment-version.sh -p $currversion)

version=${1:-$newVersion}
echo "new -> $version"

git tag $version
git push origin $version
