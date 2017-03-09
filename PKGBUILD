pkgname=(
    'vhost-git'
)
pkgver=2.0.beta.r1.g26f2a26
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
    cp -pR "${pkgname}" "${pkgdir}/usr/local/bin/vhost"
    mkdir "${pkgdir}/etc/vhost"
    cp -r "shared/*" "${pkgdir}/etc/vhost"
    cp -r "config.yaml" "${pkgdir}/etc/vhost/config.yaml.dist"
}