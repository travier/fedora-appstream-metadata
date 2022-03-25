Summary:        Operating System AppStream Metadata for Fedora Linux
Name:           fedora-appstream-metadata
Version:        37
Release:        %autorelease -p
License:        MIT
URL:            https://fedoraproject.org/
Source1:        org.fedoraproject.fedora.metainfo.xml
Source2:        update-appstream-metadata.go
BuildArch:      noarch

%description
Operating System AppStream Metadata for Fedora Linux

%prep

%build

%install
install -Dpm 0644 %{SOURCE1} %{buildroot}%{_datadir}/metainfo/org.fedoraproject.fedora.metainfo.xml

%files
%{_datadir}/metainfo/org.fedoraproject.fedora.metainfo.xml

%changelog
%autochangelog
