label = "${UUID.randomUUID().toString()}"
BUILD_FOLDER = "/go"
expired=240
quay_user = "gkirok"
quay_credentials = "iguazio-dev-quay-credentials"
docker_user = "gallziguazio"
docker_credentials = "iguazio-dev-docker-credentials"
artifactory_user = "gallz"
artifactory_url = "iguazio-dev-artifactory-url"
artifactory_credentials = "iguazio-dev-artifactory-credentials"
git_project = "logfwd"
git_project_user = "gkirok"
git_deploy_user = "iguazio-dev-git-user"
git_deploy_user_token = "iguazio-dev-git-user-token"

properties([pipelineTriggers([[$class: 'PeriodicFolderTrigger', interval: '2m']])])
podTemplate(label: "${git_project}-${label}", yaml: """
apiVersion: v1
kind: Pod
metadata:
  name: "${git_project}-${label}"
  labels:
    jenkins/kube-default: "true"
    app: "jenkins"
    component: "agent"
spec:
  shareProcessNamespace: true
  containers:
    - name: jnlp
      image: jenkins/jnlp-slave
      resources:
        limits:
          cpu: 1
          memory: 2Gi
        requests:
          cpu: 1
          memory: 2Gi
      volumeMounts:
        - name: go-shared
          mountPath: /go
    - name: docker-cmd
      image: docker
      command: [ "/bin/sh", "-c", "--" ]
      args: [ "while true; do sleep 30; done;" ]
      volumeMounts:
        - name: docker-sock
          mountPath: /var/run
        - name: go-shared
          mountPath: /go
  volumes:
    - name: docker-sock
      hostPath:
          path: /var/run
    - name: go-shared
      emptyDir: {}
"""
) {
    node("${git_project}-${label}") {
        withCredentials([
                usernamePassword(credentialsId: git_deploy_user, passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME'),
                string(credentialsId: git_deploy_user_token, variable: 'GIT_TOKEN'),
                string(credentialsId: artifactory_url, variable: 'ARTIFACTORY_URL')
        ]) {
            def TAG_VERSION

            def pipelinex = library(identifier: 'pipelinex@DEVOPS-204-pipelinex', retriever: modernSCM(
                    [$class: 'GitSCMSource',
                     credentialsId: "iguazio-dev-git-user-private-key",
                     remote: 'git@github.com:gkirok/pipelinex.git'])).com.iguazio.pipelinex

            stage('get tag data') {
                container('jnlp') {
                    TAG_VERSION = sh(
                            script: "echo ${TAG_NAME} | tr -d '\\n' | egrep '^v[\\.0-9]*.*\$' | sed 's/v//'",
                            returnStdout: true
                    ).trim()

                    sh "curl -H \"Authorization: token ${GIT_TOKEN}\" https://api.github.com/repos/${git_project_user}/${git_project}/releases/tags/v${TAG_VERSION} > ~/tag_version"

                    PUBLISHED_BEFORE = sh(
                            script: "tag_published_at=\$(cat ~/tag_version | python -c 'import json,sys;obj=json.load(sys.stdin);print obj[\"published_at\"]'); SECONDS=\$(expr \$(date +%s) - \$(date -d \"\$tag_published_at\" +%s)); expr \$SECONDS / 60 + 1",
                            returnStdout: true
                    ).trim().toInteger()

                    echo "$TAG_VERSION"
                    echo "$PUBLISHED_BEFORE"
                }
            }

            if ( TAG_VERSION != null && TAG_VERSION.length() > 0 && PUBLISHED_BEFORE < expired ) {
                stage('prepare sources') {
                    container('jnlp') {
                        sh """
                            cd ${BUILD_FOLDER}
                            git clone https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/${git_project_user}/${git_project}.git src/github.com/v3io/${git_project}
                            cd src/github.com/v3io/${git_project}
                            git checkout v${TAG_VERSION}
                        """
                    }
                }

                stage('build ${git_project} in dood') {
                    container('docker-cmd') {
                        sh """
                            cd ${BUILD_FOLDER}/src/github.com/v3io/${git_project}
                            docker build . -f Dockerfile.multi --tag ${git_project}:${TAG_VERSION} --tag ${docker_user}/${git_project}:${TAG_VERSION} --tag ${docker_user}/${git_project}:latest --tag quay.io/${quay_user}/${git_project}:${TAG_VERSION} --tag quay.io/${quay_user}/${git_project}:latest --tag ${ARTIFACTORY_URL}/${artifactory_user}/${git_project}:${TAG_VERSION} --tag ${ARTIFACTORY_URL}/${artifactory_user}/${git_project}:latest
                        """
                    }
                }

                stage('push') {
                    container('docker-cmd') {
                        dockerx.images_push_multi_registries(["${git_project}:${TAG_VERSION}"], [pipelinex.DockerRepoDev.ARTIFACTORY_K8S, pipelinex.DockerRepoDev.DOCKER_HUB, pipelinex.DockerRepoDev.QUAY_IO])
                    }
                }

                stage('update release status') {
                    sh "release_id=\$(curl -H \"Content-Type: application/json\" -H \"Authorization: token ${GIT_TOKEN}\" -X GET https://api.github.com/repos/${git_project_user}/${git_project}/releases/tags/v${TAG_VERSION} | python -c 'import json,sys;obj=json.load(sys.stdin);print obj[\"id\"]'); curl -H \"Content-Type: application/json\" -H \"Authorization: token ${GIT_TOKEN}\" -X PATCH https://api.github.com/repos/${git_project_user}/${git_project}/releases/\${release_id} -d '{\"prerelease\": false}'"
                }
            } else {
                stage('warning') {
                    if (PUBLISHED_BEFORE >= expired) {
                        echo "Tag too old, published before $PUBLISHED_BEFORE minutes."
                    } else {
                        echo "${TAG_VERSION} is not release tag."
                    }
                }
            }
        }
    }
}
