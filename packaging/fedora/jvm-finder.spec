Name:      jvm-finder
Version:   ${version}
Release:   1%{?dist}
Summary:   A tool for finding an appropriate installed JVM to run your program

License:   Apache-2.0
URL:       https://github.com/loicrouchon/jvm-finder
Source0:   https://github.com/loicrouchon/jvm-finder/archive/refs/tags/v${version}.tar.gz
# Source0:   https://github.com/loicrouchon/jvm-finder/archive/refs/heads/<BRANCH>.zip

BuildArch: x86_64 aarch64
BuildRequires: make, java-latest-openjdk-devel, golang

%description
jvm-finder is a command-line tool for finding an appropriate installed JVM to run your program.
It provides a simple and efficient way to specify what version of JVM you want
and what kind of features (java, javac, native-image, ...) it should provide.

%global debug_package %{nil}

%prep

%setup -q -n jvm-finder-${version}
%build
GO_LD_FLAGS='-linkmode=external' make test build

%install
%define distdir build
find .
mkdir -p %{buildroot}/usr/bin %{buildroot}/usr/share/%{name} %{buildroot}/usr/share/%{name}/metadata-extractor %{buildroot}/etc/%{name}
ln -s ../share/%{name}/%{name} %{buildroot}/usr/bin/%{name}
install -p -m 755 %{distdir}/go/%{name} %{buildroot}/usr/share/%{name}/%{name}
install -p -m 644 %{distdir}/classes/JvmMetadataExtractor.class %{buildroot}/usr/share/%{name}/metadata-extractor/JvmMetadataExtractor.class
install -p -m 644 packaging/fedora/config.conf %{buildroot}/etc/%{name}/config.conf

%files
%license LICENSE
/usr/bin/%{name}
/usr/share/%{name}/%{name}
/usr/share/%{name}/metadata-extractor/JvmMetadataExtractor.class
/etc/%{name}/config.conf
