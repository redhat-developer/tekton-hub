#!//bin/bash
set -e -u -o pipefail

CONFIG_JS=/usr/share/nginx/html/config.js
echo "Current ENV"
cat $CONFIG_JS
ls -l  $CONFIG_JS
echo '----------------------------------------------'

cat <<-EOF > $CONFIG_JS
window.config = {
  API_URL: '$API_URL',
  GH_CLIENT_ID: '$GH_CLIENT_ID',
};


// TODO(sthaha): remove debug message
console.debug("application config", window.config);
EOF

echo "Modified ENV"
cat $CONFIG_JS
echo '----------------------------------------------'


echo Starting Nginx
exec nginx -g 'daemon off;'
