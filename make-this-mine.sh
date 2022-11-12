if [[ "${GITHUB_WORKSPACE}" ]]; then
		echo "Repo: ${GITHUB_REPOSITORY}"
		cd ${GITHUB_WORKSPACE}
		# start replacing
#		egrep -lRZ 'Benbentwo@go-bin-generic' . | grep -v Makefile | xargs -0 -l sed -i -e 's@Benbentwo@go-bin-generic@${GITHUB_REPOSITORY}/g'
    find . -type f \( -iname \*.mod -o -iname \*.go \) -print0 | xargs -0 sed -i "s@Benbentwo/go-bin-generic@${GITHUB_REPOSITORY}@g"
		ORG=$(echo ${GITHUB_REPOSITORY} | awk -F '/' '{print $1}')
		REPO=$(echo ${GITHUB_REPOSITORY} | awk -F '/' '{print $2}')
		echo "ORG: ${ORG}, REPO: ${REPO}"
		sed -i "s@ORG			:= Benbentwo@ORG         := ${ORG}@g" Makefile
		sed -i "s@REPO        := go-bin-generic@REPO        := ${REPO}@g" Makefile
		sed -i "s@BINARY      := go-bin-generic@BINARY      := ${REPO}@g" Makefile

else  # someones local
		read -r -p "Do you want to continue? [y/N] " response
		case "$response" in
			[yY][eE][sS]|[yY])
				continue
				;;
			*)
			  echo "Exiting..."
				exit 0
				;;
		esac
		GITHUB_REPOSITORY=$(git remote get-url origin | awk -F 'github.com/' '{print $2}' | awk -F '.git' '{print $1}')
		echo "GITHUB REPOSITORY: ${GITHUB_REPOSITORY}"
		find . -type f \( -iname \*.mod -o -iname \*.go \) -print0 | xargs -0 sed -i "" "s@Benbentwo/go-bin-generic@${GITHUB_REPOSITORY}@g"
		ORG=$(echo ${GITHUB_REPOSITORY} | awk -F '/' '{print $1}')
		REPO=$(echo ${GITHUB_REPOSITORY} | awk -F '/' '{print $2}')
		echo "ORG: ${ORG}, REPO: ${REPO}"
		sed -i "" "s@ORG			:= Benbentwo@ORG			:= ${ORG}@g" Makefile
		sed -i "" "s@REPO        := go-bin-generic@REPO        := ${REPO}@g" Makefile
		sed -i "" "s@BINARY 		:= go-bin-generic@BINARY 		:= ${REPO}@g" Makefile
fi
