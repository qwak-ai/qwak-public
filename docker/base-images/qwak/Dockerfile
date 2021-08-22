FROM bentoml/model-server:latest-py38
ENV PYTHONDONTWRITEBYTECODE=true

RUN apt-get --allow-releaseinfo-change update && apt-get install --no-install-recommends -y \
  git=1:2.20.1-2+deb10u3 \
  zip=3.0-11+b1 \
  && rm -rf /var/lib/apt/lists/*
COPY environment.yml /environment.yml
RUN conda env create -f environment.yml
