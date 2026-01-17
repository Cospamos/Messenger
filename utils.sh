function nginxUpd {
  sudo nginx -t -c $(pwd)/nginx.conf
  if [ $? -eq 0 ]; then
    sudo nginx -s reload -c $(pwd)/nginx.conf
  fi
}

function nginxStart {
  sudo nginx -c $(pwd)/nginx.conf
}

function runServe {
  ~/go/bin/air
}