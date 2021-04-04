# [Choice] Go version: 1, 1.16, 1.15
ARG VARIANT=1.16
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Option] Install Node.js
ARG INSTALL_NODE="true"
ARG NODE_VERSION="lts/*"
RUN if [ "${INSTALL_NODE}" = "true" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next line to use go get to install anything else you need
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

# Desktop Setup
COPY library-scripts/desktop-lite-debian.sh /tmp/library-scripts/
RUN apt-get update && bash /tmp/library-scripts/desktop-lite-debian.sh
# Fyne Dependencies
RUN apt-get install libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev
ENV DBUS_SESSION_BUS_ADDRESS="autolaunch:" DISPLAY=":1" LANG="en_US.UTF-8" LANGUAGE="en_US.UTF-8"
ENTRYPOINT ["/usr/local/share/desktop-init.sh"]
CMD ["sleep", "infinity"]

