pipeline {
    agent any

    environment {

        // Github Container Registry Login Details...
        GHCR_USER               = credentials('gchr')

        REPO                    = "ghcr.io"
        CORE_IMAGE              = "${REPO}/ubx-training/go2rtc"

        TAG_ID                  = sh(returnStdout: true, script: "git log -1 --oneline --pretty=%h").trim()
        GIT_COMMITTER_NAME      = sh(returnStdout: true, script: "git show -s --pretty=%an").trim()
        GIT_MESSAGE             = sh(returnStdout: true, script: "git show -s --pretty=%s").trim()

        JENKINS_PROJECT         = ""
        go2rtc                  = ""

    }

    stages {


        stage("Set Variables & Configuration") {
            steps {
                script {
                    def alljob = JOB_NAME.tokenize('/') as String[]
                    JENKINS_PROJECT = alljob[0]
                }
            }
        }

        stage("Prepare Build Environment to support arm") {
            steps {
                script {
                    // Set the Docker Buildx builder here - this may change from time to time so check on the server 
                        // [root@ip-172-31-35-170 docker]$ docker buildx ls
                        // NAME/NODE           DRIVER/ENDPOINT             STATUS  BUILDKIT PLATFORMS
                        // priceless_fermat *  docker-container
                        //   priceless_fermat0 unix:///var/run/docker.sock running v0.12.5  linux/amd64*, linux/arm64*, linux/arm/v7*, linux/amd64/v2, linux/amd64/v3, linux/amd64/v4, linux/386
                        // default             docker
                        //   default           default                     running 20.10.4  linux/amd64, linux/arm64, linux/386

                    sh "docker buildx use priceless_fermat"
                }
            }
        }

        stage("Build go2rtc Docker Image") {
            parallel {
                stage('go2rtc build') {
                    steps {
                        sh 'docker buildx build --platform linux/arm64,linux/amd64,linux/arm/v7 -t "${CORE_IMAGE}:${TAG_ID}" -f Dockerfile .'
                    }
                }
            }
        }

        stage('Deployment') {
            parallel {
                stage('Dev/Stage') {
                    when {
                        branch 'master'
                    }
                    steps {
                        sh 'echo $GHCR_USER_PSW | docker login ghcr.io -u $GHCR_USER_USR --password-stdin'
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64,linux/arm/v7 -t "${CORE_IMAGE}:staging" -f Dockerfile .'
                    }
                }
                stage('Production') {
                    when {
                        branch 'production'
                    }
                    steps {
                        sh 'echo $GHCR_USER_PSW | docker login ghcr.io -u $GHCR_USER_USR --password-stdin'
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64,linux/arm/v7 -t "${CORE_IMAGE}:production" -f Dockerfile .'
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64,linux/arm/v7 -t "${CORE_IMAGE}:latest" -f Dockerfile .'
                    }
                }
            }
        }
    }

    post {
        success {
            slackSend color: "#03CC00", message: "*Build Succeeded (${env.BUILD_NUMBER})*\n Job: ${env.JOB_NAME}\n Commit: ${env.GIT_MESSAGE}\n Author: ${env.GIT_COMMITTER_NAME}\n <${env.RUN_DISPLAY_URL}|Open Jenkins Log>"
        }
        failure {
            withAWS(region:'us-east-1') {
                sh "cp ${env.JENKINS_HOME}/jobs/${JENKINS_PROJECT}/branches/${env.BRANCH_NAME}/builds/${env.BUILD_NUMBER}/log /tmp/${env.BUILD_ID}.log"
                s3Upload(
                    pathStyleAccessEnabled: true,
                    payloadSigningEnabled: true,
                    file: "/tmp/${env.BUILD_ID}.log",
                    bucket: "12rnd-ubx-build-logs",
                    // https://12rnd-ubx-build-logs.s3.amazonaws.com/jenkins/${env.JOB_NAME}/${env.BRANCH_NAME}/${env.BUILD_NUMBER}/log.txt
                    path: "jenkins/${env.JOB_NAME}/${env.BUILD_NUMBER}/log.txt",
                    acl: "PublicRead"
                )
            }
            slackSend color: "#FF0000", message: "*Build Failed (${env.BUILD_NUMBER})*\n Job: ${env.JOB_NAME}\n Commit: ${env.GIT_MESSAGE}\n Author: ${env.GIT_COMMITTER_NAME}\n <${env.RUN_DISPLAY_URL}|Open Jenkins> | <https://12rnd-ubx-build-logs.s3.amazonaws.com/jenkins/${env.JOB_NAME}/${env.BUILD_NUMBER}/log.txt|View Developer Logs>"
        }
    }
}
