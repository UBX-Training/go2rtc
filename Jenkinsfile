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
                        script {
                            go2rtc = docker.build("${CORE_IMAGE}:${TAG_ID}", "-f Dockerfile .")
                            go2rtc.push()
                        }
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
                        script {
                            // Only push images that were a successful build...
                            // We're using 'watchtower' to auto-deploy the container once it's built and pushed into the docker registry...
                            go2rtc.push("staging")
                        }
                    }
                }
                stage('Production') {
                    when {
                        branch 'production'
                    }
                    steps {
                        script {
                            go2rtc.push()
                            go2rtc.push("production")
                            go2rtc.push("latest")
                        }
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
