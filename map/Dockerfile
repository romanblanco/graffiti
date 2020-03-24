FROM ruby:2.7
# throw errors if Gemfile has been modified since Gemfile.lock
#RUN bundle config --global frozen 1
WORKDIR /usr/src/app
COPY app/Gemfile* ./
RUN bundle install
COPY app/ ./
EXPOSE 4567/tcp
CMD ["bundle", "exec", "ruby", "map.rb", "-o", "0.0.0.0"]
