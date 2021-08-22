FROM bentoml/model-server:latest-py38

# OS requirements
# hadolint ignore=DL3008
RUN apt-get --allow-releaseinfo-change update && apt-get install --no-install-recommends -y \
  git \
  zip \
  && rm -rf /var/lib/apt/lists/*

# Conda env setup
ENV PYTHONDONTWRITEBYTECODE=true
RUN conda update conda
