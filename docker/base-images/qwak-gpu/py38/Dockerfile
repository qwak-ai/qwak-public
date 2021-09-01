FROM bentoml/model-server:latest-py38-gpu

# OS requirements
# hadolint ignore=DL3008
RUN apt-get --allow-releaseinfo-change update && apt-get install --no-install-recommends -y \
  git \
  zip \
  curl \
  && rm -rf /var/lib/apt/lists/*

# Install poetry
RUN curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/get-poetry.py | python -

# Conda env setup
ENV PYTHONDONTWRITEBYTECODE=true
RUN conda update conda
RUN pip install poetry2conda
