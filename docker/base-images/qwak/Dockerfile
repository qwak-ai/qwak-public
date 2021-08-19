FROM bentoml/model-server:latest-py38
# hadolint ignore=DL3008
RUN apt-get --allow-releaseinfo-change update && apt-get install --no-install-recommends -y \
  git \
  zip \
  && rm -rf /var/lib/apt/lists/*
COPY environment.yml /environment.yml
RUN conda env create -f environment.yml && conda clean -afy
