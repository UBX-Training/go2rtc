pipeline {
    agent any

    environment {

        REPO                    = "docker-registry.aws.gymsystems.co"
        CORE_IMAGE              = "${REPO}/cctv/go2rtc"

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

        stage("Build go2rtc Docker Image") {
            parallel {
                stage('go2rtc build') {
                    steps {
                        docker-registry.aws.gymsystems.co/backoffice/render-engine:staging
                        sh 'docker buildx build --platform linux/arm64,linux/amd64 -t "${CORE_IMAGE}:${TAG_ID}" -f Dockerfile .'
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
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64 -t "${CORE_IMAGE}:staging" -f Dockerfile .'
                    }
                }
                stage('Production') {
                    when {
                        branch 'production'
                    }
                    steps {
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64 -t "${CORE_IMAGE}:production" -f Dockerfile .'
                        sh 'docker buildx build --push --platform linux/arm64,linux/amd64 -t "${CORE_IMAGE}:latest" -f Dockerfile .'
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
