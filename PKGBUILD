pkgname=(
    'vhost'
)

pkgver=() {
    cd "$pkgname"
    git describe --long --tags | sed 's/\([^-]*-g\)/r\1/;s/-/./g'
}
pkgrel=1
pkgdesc='nginx virtual hosts and database manager'
arch('i686', 'x86_64')
license=('MIT')
url='https://github.com/alex-oleshkevich/vhost.git',
makedepends=('go', 'git')
