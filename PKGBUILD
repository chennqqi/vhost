pkgname=(
    'vhost-git'
)
pkgver=2.0.beta.r3.g7737a54
pkgrel=1
pkgdesc='nginx virtual hosts and database manager'
arch=('any')
license=('MIT')
url='https://github.com/alex-oleshkevich/vhost.git',
source=("$pkgname"::'git+https://github.com/alex-oleshkevich/vhost.git')
provides=('vhost')
md5sums=('SKIP')

pkgver() {
    cd "${pkgname}"
    git describe --tags --long 2>/dev/null | sed 's/\([^-]*-g\)/r\1/;s/-/./g'
}

package() {
    cd "$srcdir/${pkgname}/"
    go get github.com/tools/godep
    godep restore
    make build
    install -Dm 644 "LICENSE" "${pkgdir}/usr/share/licenses/${pkgname}"
    install -D "${pkgname}" "${pkgdir}/usr/local/bin/vhost"
    install -D "shared/presets/default.tpl" "${pkgdir}/etc/vhost/presets/default.tpl"
    install -D "shared/presets/fpm.tpl" "${pkgdir}/etc/vhost/presets/fpm.tpl"
    install -D "shared/presets/symfony.tpl" "${pkgdir}/etc/vhost/presets/symfony.tpl"
    install -D "shared/templates/vhost.tpl" "${pkgdir}/etc/vhost/templates/vhost.tpl"
    install -D "config.yaml" "${pkgdir}/etc/vhost/config.yaml.dist"
}