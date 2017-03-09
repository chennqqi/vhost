pkgname=(
    'vhost-git'
)
pkgver=2.0.beta.r0.g3e29191
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
    pwd
    make build
    install -Dm 644 "LICENSE" "${pkgdir}/usr/share/licenses/${pkgname}"
    cp -pR "vhost" "${pkgdir}/usr/local/bin/vhost"
    mkdir "${pkgdir}/etc/vhost"
    cp -r "shared/*" "${pkgdir}/etc/vhost"
    cp -r "config.yaml" "${pkgdir}/etc/vhost/config.yaml.dist"
}